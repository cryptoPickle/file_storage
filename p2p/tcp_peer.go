package p2p

import (
	"net"
	"sync"
)

type TCPPeer struct {
	net.Conn
	// if dial outbound => true, if accept outbound => false
	outbound bool
	Wg       *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Write(b)
	return err
}

func (p *TCPPeer) CloseStream() {
	p.Wg.Done()
}
