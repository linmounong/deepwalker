package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestDumpVocab(t *testing.T) {
	vocab := map[uint32]uint64{
		1: 2,
		5: 6,
		3: 4,
	}
	buf := new(bytes.Buffer)
	out := bufio.NewWriter(buf)
	dumpVocab(vocab, out)
	s := buf.String()
	if s != "</s> 0\n5 6\n3 4\n1 2\n" {
		t.Error("wrong output", s)
	}
}
