package zenc

import (
	"io"
)

const defaultBufferSize = 4096

// Stage represents a data processor in Pipeline
type Stage interface {
	Process([]byte) []byte
}

// Pipeline represents an IO data flow
type Pipeline struct {
	ss []Stage
}

// NewPipeline creates a Pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// AddStage appends Stage s to the Pipeline
func (p *Pipeline) AddStage(s Stage) {
	p.ss = append(p.ss, s)
}

// Process makes Pipeline a Stage itself
func (p *Pipeline) Process(data []byte) []byte {
	for _, s := range p.ss {
		data = s.Process(data)
	}
	return data
}

// Run reads data from r, processes it and writes the processed data to w
func (p *Pipeline) Run(r io.Reader, w io.Writer) (n int64, err error) {
	buffer := make([]byte, defaultBufferSize)
	for {
		nr, er := r.Read(buffer)
		if nr > 0 {
			nw, ew := w.Write(p.Process(buffer[:nr]))
			if nw > 0 {
				n += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return n, err
}
