package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder    Decoder

	peerLock sync.RWMutex
	peers    map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		shakeHands: NOPHandshakeFunc,
		listenAddr: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	go t.acceptLoop()
	return nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Println("TCP Accept Error: ", err)
		}
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	if err := t.shakeHands(peer); err != nil {
	}

	msg := &Temp{}
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}
	}
	fmt.Printf("new incoming connection %+v \n", peer)
}
