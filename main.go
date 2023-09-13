package main

import (
	"fmt"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/sdp/v3"
	"github.com/pion/webrtc/v4"
)

func doSignaling(offerer, answerer *webrtc.PeerConnection) {

	offer, err := offerer.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	if err = offerer.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	if err = answerer.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	answer, err := answerer.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	if err = answerer.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	if err = offerer.SetRemoteDescription(answer); err != nil {
		panic(err)
	}
}

func main() {
	offerer, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	// Add three Simulcast layers to the Offerer
	tracks := []*webrtc.TrackLocalStaticRTP{}
	for _, id := range []string{"a", "b", "c"} {
		track, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion", webrtc.WithRTPStreamID(id))
		if err != nil {
			panic(err)
		}

		tracks = append(tracks, track)
	}

	rtpSender, err := offerer.AddTrack(tracks[0])
	if err != nil {
		panic(err)
	}

	if err = rtpSender.AddEncoding(tracks[1]); err != nil {
		panic(err)
	}

	if err = rtpSender.AddEncoding(tracks[2]); err != nil {
		panic(err)
	}

	answerer, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	offerer.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Println("PeerConnectionState", s)
	})

	answerer.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			if err = offerer.AddICECandidate(c.ToJSON()); err != nil {
				panic(err)
			}
		}
	})

	answerer.OnTrack(func(t *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		fmt.Printf("New Incoming Track RID(%s)\n", t.RID())
	})

	doSignaling(offerer, answerer)

	parameters := rtpSender.GetParameters()
	var midID, ridID uint8
	for _, extension := range parameters.HeaderExtensions {
		switch extension.URI {
		case sdp.SDESMidURI:
			midID = uint8(extension.ID)
		case sdp.SDESRTPStreamIDURI:
			ridID = uint8(extension.ID)
		}
	}

	for seqNo := uint16(0); ; seqNo++ {
		pkt := &rtp.Packet{
			Header: rtp.Header{
				Version:        2,
				SequenceNumber: seqNo,
				PayloadType:    96,
			},
			Payload: []byte{0x00, 0x00},
		}

		for i := range tracks {
			if err := pkt.SetExtension(ridID, []byte(tracks[i].RID())); err != nil {
				panic(err)
			}

			if err := pkt.SetExtension(midID, []byte("0")); err != nil {
				panic(err)
			}

			if err := tracks[i].WriteRTP(pkt); err != nil {
				panic(err)
			}
		}

		time.Sleep(time.Millisecond * 200)
	}
}
