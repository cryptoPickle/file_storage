package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/cryptoPickle/file_storage/p2p"
)

type FileServerOpts struct {
	Transport      p2p.Transporter
	StorageRoot    string
	BootstrapNodes []string
}

type FileServer struct {
	*FileServerOpts
	store *Store

	peerLock *sync.Mutex
	peers    map[string]p2p.Peer

	quitch chan struct{}
}

func defaultOptions() *FileServerOpts {
	return &FileServerOpts{
		StorageRoot: "storage",
	}
}

type FileServerOptsFn func(*FileServerOpts)

func NewFileServer(opts ...FileServerOptsFn) *FileServer {
	defaultOpts := defaultOptions()
	for _, fn := range opts {
		fn(defaultOpts)
	}

	fs := &FileServer{
		FileServerOpts: defaultOpts,
		peers:          make(map[string]p2p.Peer),
		peerLock:       &sync.Mutex{},
		quitch:         make(chan struct{}),
	}

	WithRoot := func(sopts *StoreOpts) {
		sopts.Root = defaultOpts.StorageRoot
	}

	fs.store = NewStore(WithRoot)
	return fs
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()
	s.loop()
	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[peer.RemoteAddr().String()] = peer
	log.Printf("connected with remote %s", peer.RemoteAddr())
	return nil
}

type Payload struct {
	Key  string
	Data []byte
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)
	if err := s.store.Write(key, tee); err != nil {
		return err
	}

	if _, err := io.Copy(buf, r); err != nil {
		return err
	}

	p := &Payload{
		Key:  key,
		Data: buf.Bytes(),
	}
	return s.broadcast(p)
}

func (s *FileServer) broadcast(p *Payload) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer.Conn())
	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) bootstrapNetwork() {
	for _, addr := range s.BootstrapNodes {
		go func(address string) {
			log.Println("attempting to dial: ", addr)
			if err := s.Transport.Dial(address); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}
}

func (s *FileServer) loop() {
	defer func() {
		s.Transport.Close()
		log.Println("file server stopped")
	}()
	for {
		select {
		case msg := <-s.Transport.Consume():
			var p Payload
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&p); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("recived msg %+v\n", string(p.Data))
		case <-s.quitch:
			return
		}
	}
}
