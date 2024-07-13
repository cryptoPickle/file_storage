package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

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

type MessageStoreFile struct {
	Key  string
	Size int64
}

type Message struct {
	Payload any
}

type MessageGetFile struct {
	Key string
}

func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.store.Has(key) {
		return s.store.Read(key)
	}

	fmt.Printf("don't have file (%s) locally fetching from network\n", key)
	msg := Message{
		Payload: MessageGetFile{
			Key: key,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}

	select {}
	return nil, nil
}

func (s *FileServer) Store(key string, r io.Reader) error {
	var (
		fileBuff = new(bytes.Buffer)
		tee      = io.TeeReader(r, fileBuff)
	)
	size, err := s.store.Write(key, tee)
	if err != nil {
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 50)

	for _, peer := range s.peers {
		n, err := io.Copy(peer, fileBuff)
		if err != nil {
			return err
		}
		fmt.Println("received and written bytes to disk ", n)
	}

	return nil
}

func (s *FileServer) stream(p *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(&msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
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
		log.Println("file server stopped due to error or user quit action ")
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			fmt.Println("received a message ")
			var m Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&m); err != nil {
				log.Println("decoding error: ", err)
			}

			if err := s.handleMessage(rpc.From, &m); err != nil {
				log.Println("handle message errror ", err)
			}

		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	fmt.Println("hereeeeee --------")
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, &v)
	case MessageGetFile:
		return s.handleMessageGetFile(from, &v)
	}

	return nil
}

func (s *FileServer) handleMessageGetFile(from string, msg *MessageGetFile) error {
	if !s.store.Has(msg.Key) {
		return fmt.Errorf("file (%s) does not exits on disk", msg.Key)
	}

	fmt.Printf("got file (%s) serving over the network\n", msg.Key)
	r, err := s.store.Read(msg.Key)
	if err != nil {
		return err
	}
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer %s not in map", from)
	}

	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}
	fmt.Printf("written %d bytes over the networks to %s \n", n, from)
	return nil
}

func (s *FileServer) handleMessageStoreFile(from string, msg *MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) could not be found in peer list", from)
	}
	if _, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size)); err != nil {
		return err
	}
	peer.CloseStream()
	return nil
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
