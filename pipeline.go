package zenc

import "io"

// Stage represents a data processing step of Pipeline
type Stage interface {
	io.ReadWriter
	SetNext(Stage)
	Next() Stage
}

// MiddleStage processes data and writes them to next stage
type MiddleStage struct {
	next Stage
}

// SetNext sets the next Stage
func (ms *MiddleStage) SetNext(n Stage) {
	ms.next = n
}

// Next gets the next Stage
func (ms *MiddleStage) Next() Stage {
	return ms.next
}

// Close closes the next of this Stage
func (ms *MiddleStage) Close() error {
	c, ok := ms.next.(io.Closer)
	if ok {
		return c.Close()
	}
	return nil
}

// Pipeline represents a data flow consisting of one or more stages
type Pipeline struct {
	head Stage
	tail Stage
}

// NewPipeline creates a Pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{nil, nil}
}

// AddStage adds the stage at the end of the Pipeline
func (p *Pipeline) AddStage(stage Stage) {
	if p.head == nil {
		p.head = stage
		p.tail = p.head
	} else {
		p.tail.SetNext(stage)
		p.tail = p.tail.Next()
	}
}

func (p *Pipeline) Read(buf []byte) (int, error) {
	return p.head.Read(buf)
}

func (p *Pipeline) Write(data []byte) (int, error) {
	return p.head.Write(data)
}

// Close closes the first stage of the Pipeline
func (p *Pipeline) Close() error {
	c, ok := p.head.(io.Closer)
	if ok {
		return c.Close()
	}
	return nil
}
