package zenc

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

const (
	maxChunkSize  = 4096
	chunkInfoSize = int(unsafe.Sizeof(ChunkInfo{}))
)

// ChunkInfo holds the information of a chunk
type ChunkInfo struct {
	Size int32
}

// DecodeChunkInfo converts bytes to a ChunkInfo
func DecodeChunkInfo(data []byte) ChunkInfo {
	info := ChunkInfo{}
	buf := data[:chunkInfoSize]
	info.Size = int32(binary.LittleEndian.Uint32(buf))
	return info
}

// EncodeChunkInfo converts a ChunkInfo to bytes
func EncodeChunkInfo(info ChunkInfo) []byte {
	buf := [chunkInfoSize]byte{}
	binary.LittleEndian.PutUint32(buf[:], uint32(info.Size))
	return buf[:]
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
	for cs.buffer.Len() >= maxChunkSize {
		info := ChunkInfo{maxChunkSize}
		_, err = cs.next.Write(EncodeChunkInfo(info))
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
	if cs.buffer.Len() > 0 {
		info := ChunkInfo{int32(cs.buffer.Len())}
		_, err = cs.next.Write(EncodeChunkInfo(info))
		if err != nil {
			return
		}
		_, err = cs.next.Write(cs.buffer.Bytes())
		if err != nil {
			return
		}
	}
	info := ChunkInfo{0}
	_, err = cs.next.Write(EncodeChunkInfo(info))
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
		if us.buffer.Len() < chunkInfoSize {
			break
		}
		info := DecodeChunkInfo(us.buffer.Bytes())
		chunkSize := int(info.Size)
		if chunkSize > maxChunkSize {
			err = ErrLimitExceeded
			break
		}
		if us.buffer.Len() < chunkInfoSize+chunkSize {
			break
		}
		chunk := us.buffer.Next(chunkInfoSize + chunkSize)
		_, err = us.next.Write(chunk[chunkInfoSize:])
		if err != nil {
			break
		}
	}
	return
}
