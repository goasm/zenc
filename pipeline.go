package zenc

import "io"

type Stage interface {
	io.Writer
	SetNext(Stage)
	Next() Stage
}

type Pipeline struct {
	next Stage
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) SetNext(n Stage) {
	p.next = n
}

func (p *Pipeline) Next() Stage {
	return p.next
}

func (p *Pipeline) Write(data []byte) (int, error) {
	return p.next.Write(data)
}
