package p2p

import "net"

type TCPPeer struct {
	conn net.Conn
	// if dial outbound => true, if accept outbound => false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}
