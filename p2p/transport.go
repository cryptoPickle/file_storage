package p2p

// Peer is an interface that represents remote node.
type Peer interface{}

// Transport is anything that handles the communication
// between the nodes in the network. This can be form of TCP,
// UDP, websockets ...
type Transporter interface {
	ListenAndAccept() error
}
