package main

import (
	"log"

	"github.com/cryptoPickle/file_storage/p2p"
)

type FileServerOpts struct {
	Transport   p2p.Transporter
	ListenAddr  string
	StorageRoot string
}

type FileServer struct {
	*FileServerOpts
	store *Store

	quitch chan struct{}
}

func defaultOptions() *FileServerOpts {
	return &FileServerOpts{
		ListenAddr:  ":3000",
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

	s.loop()
	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) loop() {
	defer func() {
		s.Transport.Close()
		log.Println("file server stopped")
	}()
	for {
		select {
		case msg := <-s.Transport.Consume():
			log.Println(msg)
		case <-s.quitch:
			return
		}
	}
}
