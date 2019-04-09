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
	head Stage
	tail Stage
}

// NewPipeline creates a Pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{nil, nil}
}

// AddStage adds a Stage at the end of the Pipeline
func (p *Pipeline) AddStage(stage Stage) {
	if p.head == nil {
		p.head = stage
		p.tail = p.head
	} else {
		p.tail.SetNext(stage)
	}
}

func (p *Pipeline) Write(data []byte) (int, error) {
	return p.head.Write(data)
}
