package p2p

import (
	"net"
	"sync"
)

type TCPPeer struct {
	conn net.Conn
	// if dial outbound => true, if accept outbound => false
	outbound bool
	Wg       *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *TCPPeer) Conn() net.Conn {
	return p.conn
}

func (p *TCPPeer) Read(b []byte) (int, error) {
	return p.conn.Read(b)
}

func (p *TCPPeer) CloseStream() {
	p.Wg.Done()
}
