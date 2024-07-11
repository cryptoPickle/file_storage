package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	*TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts *TCPTransportOpts) *TCPTransport {
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

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) ListenAddress() string {
	return t.ListenAddr
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Println("TCP Accept Error: ", err)
		}
		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) Dial(addr string) error {
	fmt.Print(addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	defer func() {
		fmt.Printf("dorpping peer connection %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

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

		rpc.From = conn.RemoteAddr().String()
		fmt.Println("waiting till stream  is done")
		peer.Wg.Add(1)
		t.rpcch <- rpc
		peer.Wg.Wait()
		fmt.Println("stream  done continue...")
	}
}
