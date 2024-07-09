package main

import (
	"log"

	"github.com/cryptoPickle/file_storage/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
}
