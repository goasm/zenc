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

func getSampleChunkData(limit int) *bytes.Buffer {
	result := new(bytes.Buffer)
	cs := zenc.NewChunkStage()
	cs.SetNext(zenc.NewDestStage(result))
	io.Copy(cs, io.LimitReader(bytes.NewReader(sample), int64(limit)))
	cs.Flush()
	return result
}

func putSampleChunkData(data *bytes.Buffer) *bytes.Buffer {
	result := new(bytes.Buffer)
	us := zenc.NewUnchunkStage()
	us.SetNext(zenc.NewDestStage(result))
	io.Copy(us, data)
	return result
}

func TestChunkSimple(t *testing.T) {
	total := 100
	out := getSampleChunkData(total)
	if out.Len() != total+chunkHeadSize*2 {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	size1 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size1 != 100 {
		t.Fatal("ChunkStage writes a wrong chunk length", size1)
	}
	chunk1 := out.Next(100)
	if !bytes.Equal(chunk1, sample[:100]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size2 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size2 != 0 {
		t.Fatal("ChunkStage writes a wrong chunk length", size2)
	}
}

func TestChunkLarge(t *testing.T) {
	total := 9000
	out := getSampleChunkData(total)
	if out.Len() != total+chunkHeadSize*4 {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	size1 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size1 != 4096 {
		t.Fatal("ChunkStage writes a wrong chunk length", size1)
	}
	chunk1 := out.Next(4096)
	if !bytes.Equal(chunk1, sample[:4096]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size2 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size2 != 4096 {
		t.Fatal("ChunkStage writes a wrong chunk length", size2)
	}
	chunk2 := out.Next(4096)
	if !bytes.Equal(chunk2, sample[4096:8192]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size3 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size3 != 808 {
		t.Fatal("ChunkStage writes a wrong chunk length", size3)
	}
	chunk3 := out.Next(808)
	if !bytes.Equal(chunk3, sample[8192:9000]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size4 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size4 != 0 {
		t.Fatal("ChunkStage writes a wrong chunk length", size4)
	}
}

func TestChunkFull(t *testing.T) {
	total := 8192
	out := getSampleChunkData(total)
	if out.Len() != total+chunkHeadSize*3 {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	size1 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size1 != 4096 {
		t.Fatal("ChunkStage writes a wrong chunk length", size1)
	}
	chunk1 := out.Next(4096)
	if !bytes.Equal(chunk1, sample[:4096]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size2 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size2 != 4096 {
		t.Fatal("ChunkStage writes a wrong chunk length", size2)
	}
	chunk2 := out.Next(4096)
	if !bytes.Equal(chunk2, sample[4096:8192]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size3 := int(binary.LittleEndian.Uint32(out.Next(chunkHeadSize)))
	if size3 != 0 {
		t.Fatal("ChunkStage writes a wrong chunk length", size3)
	}
}

func TestUnchunkSimple(t *testing.T) {
	total := 100
	tmp := getSampleChunkData(total)
	out := putSampleChunkData(tmp)
	if out.Len() != total {
		t.Fatal("UnchunkStage writes a wrong total length", out.Len())
	}
	chunk1 := out.Next(total)
	if !bytes.Equal(chunk1, sample[:total]) {
		t.Fatal("UnchunkStage writes mismatched chunk data")
	}
	if tmp.Len() != 0 {
		t.Fatal("UnchunkStage reads wrong length of data")
	}
}

func TestUnchunkLarge(t *testing.T) {
	total := 9000
	tmp := getSampleChunkData(total)
	out := putSampleChunkData(tmp)
	if out.Len() != total {
		t.Fatal("UnchunkStage writes a wrong total length", out.Len())
	}
	chunk1 := out.Next(total)
	if !bytes.Equal(chunk1, sample[:total]) {
		t.Fatal("UnchunkStage writes mismatched chunk data")
	}
	if tmp.Len() != 0 {
		t.Fatal("UnchunkStage reads wrong length of data")
	}
}

func TestUnchunkInput(t *testing.T) {
	total := 5000
	tmp := getSampleChunkData(total)
	end := []byte{0x00, 0x01, 0x02, 0x03}
	tmp.Write(end)
	out := putSampleChunkData(tmp)
	if out.Len() != total {
		t.Fatal("UnchunkStage writes a wrong total length", out.Len())
	}
	chunk1 := out.Next(total)
	if !bytes.Equal(chunk1, sample[:total]) {
		t.Fatal("UnchunkStage writes mismatched chunk data")
	}
	if tmp.Len() != len(end) {
		t.Fatal("UnchunkStage reads wrong length of data")
	}
	if !bytes.Equal(tmp.Bytes(), end) {
		t.Fatal("UnchunkStage reads mismatched data")
	}
}
