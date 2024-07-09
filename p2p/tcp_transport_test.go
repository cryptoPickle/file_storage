package p2p

import (
	"testing"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":8000"
	tt := NewTCPTransport(listenAddr)

	if tt.listenAddr != listenAddr {
		t.Fail()
	}

	if err := tt.ListenAndAccept(); err != nil {
		t.Fatalf("can't start listener %v", err)
	}
}
