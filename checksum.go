package zenc

import (
	"hash/crc32"
)

type checksum struct {
	sum *uint32
}

// NewChecksum creates a Stage for calculating checksum
func NewChecksum(sum *uint32) Stage {
	return &checksum{sum}
}

func (c *checksum) Process(data []byte) []byte {
	*c.sum = crc32.Update(*c.sum, crc32.IEEETable, data)
	return data
}
