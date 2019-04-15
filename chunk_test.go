package zenc_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/radonlab/zenc"
)

func fillSequence(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i)
	}
}

func TestChunkSimple(t *testing.T) {
	src := make([]byte, 100)
	fillSequence(src)
	r := bytes.NewReader(src)
	w := new(bytes.Buffer)
	cs := zenc.NewChunkStage()
	cs.SetNext(zenc.NewDestStage(w))
	io.Copy(cs, r)
	cs.Flush()
	if w.Len() != 104 {
		t.Fatal("ChunkStage writes wrong length of bytes", w.Len())
	}
	size := int(binary.LittleEndian.Uint32(w.Next(4)))
	if size != 100 {
		t.Fatal("ChunkStage writes wrong chunk size", size)
	}
	chunk := w.Next(100)
	if chunk[0] != 0 || chunk[99] != 99 {
		t.Fatal("ChunkStage writes wrong chunk data")
	}
}

func TestChunkLarge(t *testing.T) {
	src := make([]byte, 5000)
	fillSequence(src)
	r := bytes.NewReader(src)
	w := new(bytes.Buffer)
	cs := zenc.NewChunkStage()
	cs.SetNext(zenc.NewDestStage(w))
	io.Copy(cs, r)
	cs.Flush()
	if w.Len() != 5008 {
		t.Fatal("ChunkStage writes wrong length of bytes", w.Len())
	}
	size1 := int(binary.LittleEndian.Uint32(w.Next(4)))
	if size1 != 4096 {
		t.Fatal("ChunkStage writes wrong chunk size", size1)
	}
	chunk1 := w.Next(4096)
	if chunk1[0] != 0 || chunk1[4095] != 255 {
		t.Fatal("ChunkStage writes wrong chunk data")
	}
	size2 := int(binary.LittleEndian.Uint32(w.Next(4)))
	if size2 != 904 {
		t.Fatal("ChunkStage writes wrong chunk size", size2)
	}
	chunk2 := w.Next(904)
	if chunk2[0] != 0 || chunk2[903] != 135 {
		t.Fatal("ChunkStage writes wrong chunk data")
	}
}

func TestUnchunkSimple(t *testing.T) {
	prefix := [4]byte{}
	src := make([]byte, 100)
	fillSequence(src)
	binary.LittleEndian.PutUint32(prefix[:], uint32(100))
	r := new(bytes.Buffer)
	r.Write(prefix[:])
	r.Write(src)
	w := new(bytes.Buffer)
	us := zenc.NewUnchunkStage()
	us.SetNext(zenc.NewDestStage(w))
	io.Copy(us, r)
	if w.Len() != 100 {
		t.Fatal("UnchunkStage writes wrong length of bytes", w.Len())
	}
	chunk := w.Next(100)
	if chunk[0] != 0 || chunk[99] != 99 {
		t.Fatal("UnchunkStage writes wrong chunk data")
	}
}
