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
func (p *Pipeline) Run(r io.Reader, w io.Writer) {
	buffer := make([]byte, defaultBufferSize)
	for {
		nr, err := r.Read(buffer)
		if err == io.EOF {
			break
		}
		nw, err := w.Write(p.Process(buffer[:nr]))
		if nw != nr || err != nil {
			// TODO:(error handle)
			panic(err)
		}
	}
}
