package zenc

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	maxChunkSize = 4096
)

type chunkReader struct {
	b bytes.Buffer
	t []byte
	r io.Reader
}

// NewChunkReader creates a Reader that reads chunks into continuous bytes
func NewChunkReader(r io.Reader) io.Reader {
	return &chunkReader{bytes.Buffer{}, make([]byte, maxChunkSize), r}
}

func (c *chunkReader) Read(dst []byte) (n int, err error) {
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

type chunkWriter struct {
	b bytes.Buffer
	w io.Writer
}

// NewChunkWriter creates a Writer that writes continuous bytes to chunks
func NewChunkWriter(w io.Writer) io.Writer {
	return &chunkWriter{bytes.Buffer{}, w}
}

func (c *chunkWriter) Write(src []byte) (n int, err error) {
	n, err = c.b.Write(src)
	if err != nil {
		return
	}
	prefix := [4]byte{}
	for c.b.Len() >= maxChunkSize {
		binary.LittleEndian.PutUint32(prefix[:], uint32(maxChunkSize))
		_, err = c.w.Write(prefix[:])
		if err != nil {
			break
		}
		_, err = c.w.Write(c.b.Next(maxChunkSize))
		if err != nil {
			break
		}
	}
	return
}

func (c *chunkWriter) Flush() (err error) {
	prefix := [4]byte{}
	binary.LittleEndian.PutUint32(prefix[:], uint32(c.b.Len()))
	_, err = c.w.Write(prefix[:])
	if err != nil {
		return
	}
	_, err = c.w.Write(c.b.Bytes())
	if err != nil {
		return
	}
	return
}
