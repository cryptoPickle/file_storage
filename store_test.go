package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

var defaultFolder = "picklenetwork"

func TestPathTransformFunc(t *testing.T) {
	key := "somepiture.jpg"
	pathname := DefaultTransformFunc(key, defaultFolder)
	expectedPathName := defaultFolder + "/1fc6c/388d4/f5030/0ab78/9cc66/2adf1/b9db4/97b45"
	expectedOriginalKey := "1fc6c388d4f50300ab789cc662adf1b9db497b45"

	if pathname.PathName != expectedPathName {
		t.Errorf("expected %v got %v", expectedPathName, pathname.PathName)
	}

	if pathname.FileName != expectedOriginalKey {
		t.Errorf("expected %v got %v", expectedOriginalKey, pathname.FileName)
	}
}

func TestStore(t *testing.T) {
	s := NewStore()
	key := "somepicture"
	defer teardown(t, s)

	data := []byte("somedata")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	fmt.Print(string(b))
	if string(b) != string(data) {
		t.Errorf("want %s have %s ", data, b)
	}
}

func TestHas(t *testing.T) {
	s := NewStore()
	key := "somepicture"

	defer teardown(t, s)

	data := []byte("somedata")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Error("has should return true")
	}

	if ok := s.Has("not exits"); ok {
		t.Error("has should return false")
	}
}

func TestDelete(t *testing.T) {
	s := NewStore()
	key := "pictureKey"

	defer teardown(t, s)

	for i := 0; i <= 50; i++ {
		buf := make([]byte, 10*i)

		if _, err := rand.Read(buf); err != nil {
			t.Error(err)
		}

		if err := s.writeStream(key, bytes.NewReader(buf)); err != nil {
			t.Error(err)
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expected to not the have the key")
		}

	}
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
