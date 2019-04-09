package zenc

import "io"

// Stage represents a data processing step of a Pipeline
type Stage interface {
	io.Writer
	SetNext(Stage)
	Next() Stage
}

// Pipeline represents a data flow consisting of one or more Stages
type Pipeline struct {
	next Stage
}

// NewPipeline creates a Pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// SetNext sets the next Stage
func (p *Pipeline) SetNext(n Stage) {
	p.next = n
}

// Next gets the next Stage
func (p *Pipeline) Next() Stage {
	return p.next
}

func (p *Pipeline) Write(data []byte) (int, error) {
	return p.next.Write(data)
}
