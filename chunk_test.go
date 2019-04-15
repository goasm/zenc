package zenc_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/radonlab/zenc"
)

const (
	chunkHeadSize int = 4
)

var (
	sample []byte
)

func init() {
	sample = make([]byte, 10000)
	fillSequence(sample)
}

func fillSequence(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i)
	}
}

func getTestChunk(limit int) *bytes.Buffer {
	data := new(bytes.Buffer)
	cs := zenc.NewChunkStage()
	cs.SetNext(zenc.NewDestStage(data))
	io.Copy(cs, io.LimitReader(bytes.NewReader(sample), int64(limit)))
	cs.Flush()
	return data
}

func TestChunkSimple(t *testing.T) {
	total := 100
	out := getTestChunk(total)
	if out.Len() != total+chunkHeadSize {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	size1 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size1 != 100 {
		t.Fatal("ChunkStage writes a wrong chunk length", size1)
	}
	chunk1 := out.Next(100)
	if chunk1[0] != 0 || chunk1[99] != 99 {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
}

func TestChunkLarge(t *testing.T) {
	total := 5000
	out := getTestChunk(total)
	if out.Len() != total+chunkHeadSize*2 {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	size1 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size1 != 4096 {
		t.Fatal("ChunkStage writes a wrong chunk length", size1)
	}
	chunk1 := out.Next(4096)
	if chunk1[0] != 0 || chunk1[4095] != 255 {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size2 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size2 != 904 {
		t.Fatal("ChunkStage writes a wrong chunk length", size2)
	}
	chunk2 := out.Next(904)
	if chunk2[0] != 0 || chunk2[903] != 135 {
		t.Fatal("ChunkStage writes mismatched chunk data")
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
		t.Fatal("UnchunkStage writes a wrong total length", w.Len())
	}
	chunk := w.Next(100)
	if chunk[0] != 0 || chunk[99] != 99 {
		t.Fatal("UnchunkStage writes mismatched chunk data")
	}
}
