package main

import (
	"log"

	"github.com/cryptoPickle/file_storage/p2p"
)

func main() {
	s := NewFileServer(WithListenAddr, WithStorageRoot, WithTransporter)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	select {}
}

func WithListenAddr(opts *FileServerOpts) {
	opts.ListenAddr = ":1337"
}

func WithStorageRoot(opts *FileServerOpts) {
	opts.StorageRoot = "pickle-storage"
}

func WithTransporter(opts *FileServerOpts) {
	tpopts := p2p.TCPTransportOpts{
		ListenAddr:    opts.ListenAddr,
		Decoder:       p2p.NOPDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
	}

	tp := p2p.NewTCPTransport(tpopts)
	opts.Transport = tp
}
