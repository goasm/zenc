package zenc

import "errors"

// zenc errors
var (
	ErrWrongChecksum = errors.New("wrong checksum")
	ErrLimitExceeded = errors.New("limit exceeded")
)
