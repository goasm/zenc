package zenc

import "io"

// DestStage outputs data to a destination Writer
type DestStage struct {
	dest io.Writer
}

// NewDestStage creates a DestStage from a Writer
func NewDestStage(d io.Writer) *DestStage {
	return &DestStage{d}
}

// SetNext on DestStage will panic
func (ds *DestStage) SetNext(n Stage) {
	panic("zenc: cannot set next on DestStage")
}

// Next of DestStage is always nil
func (ds *DestStage) Next() Stage {
	return nil
}

func (ds *DestStage) Write(data []byte) (int, error) {
	return ds.dest.Write(data)
}
