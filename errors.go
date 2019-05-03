package zenc

import "errors"

// zenc errors
var (
	ErrWrongFileType = errors.New("wrong file type")
	ErrWrongChecksum = errors.New("wrong checksum")
	ErrLimitExceeded = errors.New("limit exceeded")
)
