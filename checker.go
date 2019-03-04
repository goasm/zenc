package zenc

import "hash/crc32"

const polynomial uint32 = 0xEDB88320

type checker struct {
	sum uint32
	tab *crc32.Table
}

// NewChecker creates a Stage for calculating checksum
func NewChecker() Stage {
	return &checker{0, crc32.MakeTable(polynomial)}
}

func (c *checker) Process(data []byte) []byte {
	c.sum = crc32.Update(c.sum, c.tab, data)
	return data
}

func (c *checker) Sum() uint32 {
	return c.sum
}
