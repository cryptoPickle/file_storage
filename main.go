package main

import (
	"fmt"
	"log"

	"github.com/cryptoPickle/file_storage/p2p"
)

func OnPeer(peer p2p.Peer) error {
	// peer.Close()
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()
	select {}
}
