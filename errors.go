package zenc

import "errors"

var ErrWrongChecksum = errors.New("wrong checksum")
var ErrLimitExceeded = errors.New("limit exceeded")
