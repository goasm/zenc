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
	r io.Reader
}

// NewChunkReader creates a Reader that reads chunks into continuous bytes
func NewChunkReader(r io.Reader) io.Reader {
	return &chunkReader{r}
}

func (c *chunkReader) Read(dst []byte) (n int, err error) {
	n, err = c.r.Read(dst)
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
	for c.b.Len() < maxChunkSize {
		_, err = c.b.Write(src)
		if err != nil {
			break
		}
	}
	nw := 0
	prefix := [4]byte{}
	chunk := []byte{}
	for c.b.Len() > 0 {
		chunk = c.b.Next(maxChunkSize)
		binary.LittleEndian.PutUint32(prefix[:], uint32(len(chunk)))
		nw, err = c.w.Write(prefix[:])
		n += nw
		if err != nil {
			break
		}
		nw, err = c.w.Write(chunk)
		n += nw
		if err != nil {
			break
		}
	}
	return
}
