package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	*StoreOpts
}

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

func NewStore(options ...OptionsFn) *Store {
	opts := defaultOptios()

	for _, fn := range options {
		fn(opts)
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	ok := s.Has(key)
	if ok {
		return os.RemoveAll(pathKey.FullPath())
	}
	defer func() {
		if ok {
			log.Printf("deleted [%s] from disk", pathKey.FileName)
		}
	}()

	return errors.New("can't delete the file")
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)

	_, err := os.Stat(pathKey.FullPath())

	return err == fs.ErrNotExist
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	filePath := pathKey.FullPath()
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk %s \n", n, filePath)
	return nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(pathKey.FullPath())
}

type OptionsFn func(*StoreOpts)

func defaultOptios() *StoreOpts {
	return &StoreOpts{
		PathTransformFunc: DefaultTransformFunc,
	}
}

type PathTransformFunc func(string) PathKey

func DefaultTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blocksize := 5
	sliceLenght := len(hashStr) / blocksize
	paths := make([]string, sliceLenght)

	for i := 0; i < sliceLenght; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}
