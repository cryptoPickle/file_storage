package p2p

import (
	"testing"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr:    ":3000",
		Decoder:       GOBDecoder{},
		HandshakeFunc: NOPHandshakeFunc,
	}
	tt := NewTCPTransport(opts)

	if tt.ListenAddr != ":3000" {
		t.Fail()
	}

	if err := tt.ListenAndAccept(); err != nil {
		t.Fatalf("can't start listener %v", err)
	}
}
