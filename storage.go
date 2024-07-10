package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "picklenetwork"

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
	Root              string
}

type Store struct {
	*StoreOpts
}

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) TopParent() (string, error) {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return "", errors.New("not valid path")
	}

	return fmt.Sprintf("%v/%v", paths[0], paths[1]), nil
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
	pathKey := s.PathTransformFunc(key, s.Root)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	path, err := pathKey.TopParent()
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key, s.Root)

	_, err := os.Stat(pathKey.FullPath())

	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
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
	pathKey := s.PathTransformFunc(key, s.Root)
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
	pathKey := s.PathTransformFunc(key, s.Root)
	return os.Open(pathKey.FullPath())
}

type OptionsFn func(*StoreOpts)

func defaultOptios() *StoreOpts {
	return &StoreOpts{
		PathTransformFunc: DefaultTransformFunc,
		Root:              defaultRootFolderName,
	}
}

type PathTransformFunc func(string, string) PathKey

func DefaultTransformFunc(key string, root string) PathKey {
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
		PathName: root + "/" + strings.Join(paths, "/"),
		FileName: hashStr,
	}
}
