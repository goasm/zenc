package zenc

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	maxChunkSize = 4096
)

type ChunkStage struct {
	MiddleStage
	buffer bytes.Buffer
}

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

type Reader struct {
	b bytes.Buffer
	t []byte
	r io.Reader
}

// NewChunkReader creates a Reader that reads chunks into continuous bytes
func NewChunkReader(r io.Reader) *Reader {
	return &Reader{bytes.Buffer{}, make([]byte, maxChunkSize), r}
}

func (c *Reader) Read(dst []byte) (n int, err error) {
	m := len(dst)
	prefix := [4]byte{}
	for c.b.Len() < m {
		_, err = c.r.Read(prefix[:])
		if err != nil {
			break
		}
		l := int64(binary.LittleEndian.Uint32(prefix[:]))
		_, err = io.CopyBuffer(&c.b, io.LimitReader(c.r, l), c.t)
		if err != nil {
			break
		}
	}
	n, err = c.b.Read(dst)
	return
}
