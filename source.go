package zenc

import "io"

// SourceStage gets input data from a source Reader
type SourceStage struct {
	src io.Reader
}

// NewSourceStage creates a SourceStage from a Reader
func NewSourceStage(s io.Reader) *SourceStage {
	return &SourceStage{s}
}

// SetNext on SourceStage will panic
func (ss *SourceStage) SetNext(n Stage) {
	panic("zenc: cannot set next on SourceStage")
}

// Next of SourceStage is always nil
func (ss *SourceStage) Next() Stage {
	return nil
}

func (ss *SourceStage) Read(buf []byte) (int, error) {
	return ss.src.Read(buf)
}

func (ss *SourceStage) Write(data []byte) (int, error) {
	panic("zenc: cannot write to SourceStage")
}
