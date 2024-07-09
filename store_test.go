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

	if pathname.Pathname != expectedPathName {
		t.Errorf("expected %v got %v", expectedPathName, pathname.Pathname)
	}

	if pathname.Original != expectedOriginalKey {
		t.Errorf("expected %v got %v", expectedOriginalKey, pathname.Original)
	}
}

func TestStore(t *testing.T) {
	s := NewStore()
	key := "somepicture"

	data := bytes.NewReader([]byte("somedata"))
	p := DefaultTransformFunc(key)
	if err := s.writeStream(key, data); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll(p.Pathname); err != nil {
		t.Errorf("can't clean the test %s", err)
	}
}
