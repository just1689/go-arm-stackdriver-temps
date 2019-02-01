package main

import (
	"bytes"
	"testing"
)

func TestReadFile(t *testing.T) {
	b, err := readFile("temp")
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(b, []byte("49000")) != 0 {
		t.Error("Did not read file")
	}
}
