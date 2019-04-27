package zenc_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
	"unsafe"

	"github.com/radonlab/zenc"
)

const chunkInfoSize = int(unsafe.Sizeof(zenc.ChunkInfo{}))

var rawData []byte

func init() {
	rawData = make([]byte, 10000)
	fillSequence(rawData)
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
	io.Copy(cs, io.LimitReader(bytes.NewReader(rawData), int64(limit)))
	cs.Close()
	return result
}

func putSampleChunkData(data *bytes.Buffer) *bytes.Buffer {
	result := new(bytes.Buffer)
	cs := zenc.NewChunkStage()
	cs.SetNext(zenc.NewSourceStage(data))
	io.Copy(result, cs)
	cs.Close()
	return result
}

func TestWriteSingleChunk(t *testing.T) {
	total := 100
	out := getSampleChunkData(total)
	if out.Len() != total+chunkInfoSize*2 {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	size1 := int(binary.LittleEndian.Uint32(out.Next(chunkInfoSize)))
	if size1 != 100 {
		t.Fatal("ChunkStage writes a wrong chunk length", size1)
	}
	chunk1 := out.Next(100)
	if !bytes.Equal(chunk1, rawData[:100]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size2 := int(binary.LittleEndian.Uint32(out.Next(chunkInfoSize)))
	if size2 != 0 {
		t.Fatal("ChunkStage writes a wrong chunk length", size2)
	}
}

func TestWriteMultipleChunks(t *testing.T) {
	total := 9000
	out := getSampleChunkData(total)
	if out.Len() != total+chunkInfoSize*4 {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	size1 := int(binary.LittleEndian.Uint32(out.Next(chunkInfoSize)))
	if size1 != 4096 {
		t.Fatal("ChunkStage writes a wrong chunk length", size1)
	}
	chunk1 := out.Next(4096)
	if !bytes.Equal(chunk1, rawData[:4096]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size2 := int(binary.LittleEndian.Uint32(out.Next(chunkInfoSize)))
	if size2 != 4096 {
		t.Fatal("ChunkStage writes a wrong chunk length", size2)
	}
	chunk2 := out.Next(4096)
	if !bytes.Equal(chunk2, rawData[4096:8192]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size3 := int(binary.LittleEndian.Uint32(out.Next(chunkInfoSize)))
	if size3 != 808 {
		t.Fatal("ChunkStage writes a wrong chunk length", size3)
	}
	chunk3 := out.Next(808)
	if !bytes.Equal(chunk3, rawData[8192:9000]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size4 := int(binary.LittleEndian.Uint32(out.Next(chunkInfoSize)))
	if size4 != 0 {
		t.Fatal("ChunkStage writes a wrong chunk length", size4)
	}
}

func TestWriteFullChunk(t *testing.T) {
	total := 8192
	out := getSampleChunkData(total)
	if out.Len() != total+chunkInfoSize*3 {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	size1 := int(binary.LittleEndian.Uint32(out.Next(chunkInfoSize)))
	if size1 != 4096 {
		t.Fatal("ChunkStage writes a wrong chunk length", size1)
	}
	chunk1 := out.Next(4096)
	if !bytes.Equal(chunk1, rawData[:4096]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size2 := int(binary.LittleEndian.Uint32(out.Next(chunkInfoSize)))
	if size2 != 4096 {
		t.Fatal("ChunkStage writes a wrong chunk length", size2)
	}
	chunk2 := out.Next(4096)
	if !bytes.Equal(chunk2, rawData[4096:8192]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	size3 := int(binary.LittleEndian.Uint32(out.Next(chunkInfoSize)))
	if size3 != 0 {
		t.Fatal("ChunkStage writes a wrong chunk length", size3)
	}
}

func TestReadSingleChunk(t *testing.T) {
	total := 100
	tmp := getSampleChunkData(total)
	out := putSampleChunkData(tmp)
	if out.Len() != total {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	chunk1 := out.Next(total)
	if !bytes.Equal(chunk1, rawData[:total]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	if tmp.Len() != 0 {
		t.Fatal("ChunkStage reads wrong length of data")
	}
}

func TestReadMultipleChunks(t *testing.T) {
	total := 9000
	tmp := getSampleChunkData(total)
	out := putSampleChunkData(tmp)
	if out.Len() != total {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	chunk1 := out.Next(total)
	if !bytes.Equal(chunk1, rawData[:total]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	if tmp.Len() != 0 {
		t.Fatal("ChunkStage reads wrong length of data")
	}
}

func TestReadInputSize(t *testing.T) {
	total := 5000
	tmp := getSampleChunkData(total)
	end := []byte{0x00, 0x01, 0x02, 0x03}
	tmp.Write(end)
	out := putSampleChunkData(tmp)
	if out.Len() != total {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	chunk1 := out.Next(total)
	if !bytes.Equal(chunk1, rawData[:total]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	if tmp.Len() != len(end) {
		t.Fatal("ChunkStage reads wrong length of data")
	}
	if !bytes.Equal(tmp.Bytes(), end) {
		t.Fatal("ChunkStage reads mismatched data")
	}
}

func TestReadIncompleteInput(t *testing.T) {
	total := 100
	tmp := getSampleChunkData(total)
	tmp.Truncate(tmp.Len() - chunkInfoSize)
	out := putSampleChunkData(tmp)
	if out.Len() != total {
		t.Fatal("ChunkStage writes a wrong total length", out.Len())
	}
	chunk1 := out.Next(total)
	if !bytes.Equal(chunk1, rawData[:total]) {
		t.Fatal("ChunkStage writes mismatched chunk data")
	}
	if tmp.Len() != 0 {
		t.Fatal("ChunkStage reads wrong length of data")
	}
}
