package main

import (
	"bytes"
	"os"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "somepiture.jpg"
	pathname := DefaultTransformFunc(key)
	expectedPathName := "1fc6c/388d4/f5030/0ab78/9cc66/2adf1/b9db4/97b45"
	expectedOriginalKey := "1fc6c388d4f50300ab789cc662adf1b9db497b45"

	if pathname.PathName != expectedPathName {
		t.Errorf("expected %v got %v", expectedPathName, pathname.PathName)
	}

	if pathname.FileName != expectedOriginalKey {
		t.Errorf("expected %v got %v", expectedOriginalKey, pathname.FileName)
	}
}

func RemoveFiles(t *testing.T, path string) {
	if err := os.RemoveAll(path); err != nil {
		t.Errorf("can't clean the test %s", err)
	}
}

// func TestStore(t *testing.T) {
// 	s := NewStore()
// 	key := "somepicture"
// 	p := DefaultTransformFunc(key)
// 	defer RemoveFiles(t, p.PathName)
//
// 	data := []byte("somedata")
//
// 	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
// 		t.Error(err)
// 	}
//
// 	r, err := s.Read(key)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	b, _ := io.ReadAll(r)
//
// 	fmt.Print(string(b))
// 	if string(b) != string(data) {
// 		t.Errorf("want %s have %s ", data, b)
// 	}
// }

func TestDelete(t *testing.T) {
	s := NewStore()
	key := "somepicture"

	data := []byte("somedata")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}
