package zenc

import (
	"io"
)

type chunkReader struct {
}

func NewChunkReader(r io.Reader) io.Reader {
	return &chunkReader{}
}

func (c *chunkReader) Read(dest []byte) (n int, err error) {
	return
}

type chunkWriter struct {
}

func NewChunkWriter(w io.Writer) io.Writer {
	return &chunkWriter{}
}

func (c *chunkWriter) Write(src []byte) (n int, err error) {
	return
}
