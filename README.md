# pion-bidirectional-simulcast

This example demonstrates how Pion WebRTC can both send and receieve Simulcast

### Running

* `go install github.com/sean-der/pion-bidirectional-simulcast@latest`
* `~/go/bin/pion-bidirectional-simulcast`

In the command line you should see

```
PeerConnectionState connecting
PeerConnectionState connected
New Incoming Track RID(a)
New Incoming Track RID(c)
New Incoming Track RID(b)
```

This means that the PeerConnection has started and connected succesfully. After it connected it then started sending and receiving two distinct Simulcast tracks.

