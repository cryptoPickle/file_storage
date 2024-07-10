package main

import (
	"bytes"
	"log"
	"time"

	"github.com/cryptoPickle/file_storage/p2p"
)

func main() {
	s1 := makeServer(":3000")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()

	go s2.Start()

	time.Sleep(1 * time.Second)
	data := bytes.NewReader([]byte("some private data"))
	s2.StoreData("somekey", data)
	select {}
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	WithStorageRoot := func(opts *FileServerOpts) {
		opts.StorageRoot = listenAddr + "_network"
	}

	tpopts := &p2p.TCPTransportOpts{
		Decoder:       p2p.NOPDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
		ListenAddr:    listenAddr,
	}
	WithTransporter := func(opts *FileServerOpts) {
		tp := p2p.NewTCPTransport(tpopts)
		opts.Transport = tp
	}

	WithBootstrapNodes := func(opts *FileServerOpts) {
		opts.BootstrapNodes = nodes
	}
	s := NewFileServer(
		WithStorageRoot,
		WithTransporter,
		WithBootstrapNodes,
	)

	tpopts.OnPeer = s.OnPeer

	return s
}
