package zenc

import (
	"bytes"
	"encoding/binary"
	"io"
)

const maxChunkSize = 4096

// ChunkInfo holds the information of a chunk
type ChunkInfo struct {
	Size int32
}

// ReadFrom reads the data of ChunkInfo from r
func (ci *ChunkInfo) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.LittleEndian, ci)
	if err == nil {
		n = int64(binary.Size(ci))
	}
	return
}

// WriteTo writes the data of ChunkInfo to w
func (ci *ChunkInfo) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.LittleEndian, ci)
	if err == nil {
		n = int64(binary.Size(ci))
	}
	return
}

// ChunkStage encodes contiguous bytes into data chunks
type ChunkStage struct {
	MiddleStage
	buffer bytes.Buffer
	bucket []byte
	ended  bool
}

// NewChunkStage creates a Stage for encoding data
func NewChunkStage() *ChunkStage {
	return &ChunkStage{MiddleStage{}, bytes.Buffer{}, nil, false}
}

func (cs *ChunkStage) Read(buf []byte) (n int, err error) {
	if cs.ended {
		err = io.EOF
		return
	}
	for cs.buffer.Len() < len(buf) {
		info := ChunkInfo{}
		_, er := info.ReadFrom(cs.Next())
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		if info.Size == 0 {
			// until chunk terminator reached
			cs.ended = true
			break
		} else if info.Size > maxChunkSize {
			// max chunk size exceeded
			err = ErrLimitExceeded
			break
		}
		if cs.bucket == nil {
			cs.bucket = make([]byte, maxChunkSize)
		}
		chunk := cs.bucket[:info.Size]
		_, err = cs.Next().Read(chunk)
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
		_, err = info.WriteTo(cs.Next())
		if err != nil {
			break
		}
		_, err = cs.Next().Write(cs.buffer.Next(maxChunkSize))
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
		_, err = info.WriteTo(cs.Next())
		if err != nil {
			return
		}
		_, err = cs.Next().Write(cs.buffer.Bytes())
		if err != nil {
			return
		}
	}
	// writes chunk terminator
	info := ChunkInfo{0}
	_, err = info.WriteTo(cs.Next())
	return
}
