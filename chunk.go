package zenc

import (
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
	w io.Writer
}

// NewChunkWriter creates a Writer that writes continuous bytes to chunks
func NewChunkWriter(w io.Writer) io.Writer {
	return &chunkWriter{w}
}

func (c *chunkWriter) Write(src []byte) (n int, err error) {
	n, err = c.w.Write(src)
	return
}
