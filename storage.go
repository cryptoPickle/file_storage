package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
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
	Pathname string
	Original string
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

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil {
		return err
	}

	filePath := pathKey.Filename()
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
		Pathname: strings.Join(paths, "/"),
		Original: hashStr,
	}
}

func (p PathKey) Filename() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Original)
}
