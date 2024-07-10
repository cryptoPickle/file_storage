package main

import "github.com/cryptoPickle/file_storage/p2p"

type FileServerOpts struct {
	Transport   p2p.Transporter
	ListenAddr  string
	StorageRoot string
}

type FileServer struct {
	*FileServerOpts
	store *Store
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
	}

	WithRoot := func(sopts *StoreOpts) {
		sopts.Root = defaultOpts.StorageRoot
	}

	fs.store = NewStore(WithRoot)
	return fs
}

func (s *FileServerOpts) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}
