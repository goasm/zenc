package zenc

import (
	"bytes"
	"encoding/binary"
	"io"
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
	ended  bool
}

// NewChunkStage creates a Stage for encoding data
func NewChunkStage() *ChunkStage {
	return &ChunkStage{MiddleStage{}, bytes.Buffer{}, false}
}

func (cs *ChunkStage) Read(buf []byte) (n int, err error) {
	if cs.ended {
		err = io.EOF
		return
	}
	for cs.buffer.Len() < len(buf) {
		prefix := [chunkInfoSize]byte{}
		_, er := cs.next.Read(prefix[:])
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		info := DecodeChunkInfo(prefix[:])
		if info.Size == 0 {
			cs.ended = true
			break
		}
		chunk := make([]byte, info.Size)
		_, err = cs.next.Read(chunk)
		if err != nil {
			break
		}
		_, err = cs.buffer.Write(chunk)
		if err != nil {
			break
		}
	}
	if err != nil {
		return
	}
	n, err = cs.buffer.Read(buf)
	return
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
