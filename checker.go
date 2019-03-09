package zenc

import (
	"hash/crc32"
)

type checker struct {
	sum uint32
}

// NewChecker creates a Stage for calculating checksum
func NewChecker() Stage {
	return &checker{0}
}

func (c *checker) Process(data []byte) []byte {
	c.sum = crc32.Update(c.sum, crc32.IEEETable, data)
	return data
}

func (c *checker) Sum() uint32 {
	return c.sum
}
