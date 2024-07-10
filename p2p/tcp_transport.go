package p2p

import (
	"fmt"
	"log"
	"net"
	"reflect"
)

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	log.Printf("TCP transport listening %v\n", t.ListenAddr)

	go t.acceptLoop()
	return nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Println("TCP Accept Error: ", err)
		}
		fmt.Printf("new incoming connection %+v\n", conn)
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("dorpping peer connection %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)

		if _, ok := err.(*net.OpError); ok {
			fmt.Printf("connection closed by peer \n")
			return
		}

		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			fmt.Printf("TCP read error: %s\n", err)
			fmt.Println(reflect.TypeOf(err))
			continue

		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}
}
