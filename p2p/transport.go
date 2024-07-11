package p2p

import "net"

// Peer is an interface that represents remote node.
type Peer interface {
	Send([]byte) error
	RemoteAddr() net.Addr
	Close() error
	Conn() net.Conn
	Read([]byte) (int, error)
	CloseStream()
}

// Transport is anything that handles the communication
// between the nodes in the network. This can be form of TCP,
// UDP, websockets ...
type Transporter interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
	ListenAddress() string
}
