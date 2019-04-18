package zenc

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	maxChunkSize = 4096
)

// ChunkInfo holds the information of a chunk
type ChunkInfo struct {
	Size int32
}

// DecodeChunkInfo converts bytes to a ChunkInfo
func DecodeChunkInfo(data []byte) ChunkInfo {
	info := ChunkInfo{}
	buf := data[:binary.Size(info)]
	info.Size = int32(binary.LittleEndian.Uint32(buf))
	return info
}

// EncodeChunkInfo converts a ChunkInfo to bytes
func EncodeChunkInfo(info ChunkInfo) []byte {
	buf := make([]byte, binary.Size(info))
	binary.LittleEndian.PutUint32(buf, uint32(info.Size))
	return buf
}

// ChunkStage encodes contiguous bytes into data chunks
type ChunkStage struct {
	MiddleStage
	buffer bytes.Buffer
}

// NewChunkStage creates a Stage for encoding data
func NewChunkStage() *ChunkStage {
	return &ChunkStage{MiddleStage{}, bytes.Buffer{}}
}

func (cs *ChunkStage) Write(data []byte) (n int, err error) {
	n, err = cs.buffer.Write(data)
	if err != nil {
		return
	}
	prefix := [4]byte{}
	for cs.buffer.Len() >= maxChunkSize {
		binary.LittleEndian.PutUint32(prefix[:], uint32(maxChunkSize))
		_, err = cs.next.Write(prefix[:])
		if err != nil {
			break
		}
		_, err = cs.next.Write(cs.buffer.Next(maxChunkSize))
		if err != nil {
			break
		}
	}
	return
}

// Flush writes the rest of data to the underlying writer
func (cs *ChunkStage) Flush() (err error) {
	prefix := [4]byte{}
	binary.LittleEndian.PutUint32(prefix[:], uint32(cs.buffer.Len()))
	_, err = cs.next.Write(prefix[:])
	if err != nil {
		return
	}
	_, err = cs.next.Write(cs.buffer.Bytes())
	if err != nil {
		return
	}
	return
}

// UnchunkStage decodes data chunks to contiguous bytes
type UnchunkStage struct {
	MiddleStage
	buffer bytes.Buffer
}

// NewUnchunkStage creates a Stage for decoding data
func NewUnchunkStage() *UnchunkStage {
	return &UnchunkStage{MiddleStage{}, bytes.Buffer{}}
}

func (us *UnchunkStage) Write(data []byte) (n int, err error) {
	n, err = us.buffer.Write(data)
	if err != nil {
		return
	}
	for {
		if us.buffer.Len() < 4 {
			break
		}
		prefix := us.buffer.Bytes()[:4]
		chunkSize := int(binary.LittleEndian.Uint32(prefix))
		if chunkSize > maxChunkSize {
			err = io.ErrShortBuffer
			break
		}
		if us.buffer.Len() < 4+chunkSize {
			break
		}
		chunk := us.buffer.Next(4 + chunkSize)
		_, err = us.next.Write(chunk[4:])
		if err != nil {
			break
		}
	}
	return
}
