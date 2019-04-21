package zenc

import "hash/crc32"

// ChecksumStage maintains a checksum of the data being processed
type ChecksumStage struct {
	MiddleStage
	Sum uint32
}

// NewChecksumStage creates a Stage for calculating checksum
func NewChecksumStage() *ChecksumStage {
	return &ChecksumStage{MiddleStage{}, 0}
}

func (cs *ChecksumStage) Read(buf []byte) (int, error) {
	n, err := cs.next.Read(buf)
	cs.Sum = crc32.Update(cs.Sum, crc32.IEEETable, buf[:n])
	return n, err
}

func (cs *ChecksumStage) Write(data []byte) (int, error) {
	cs.Sum = crc32.Update(cs.Sum, crc32.IEEETable, data)
	return cs.next.Write(data)
}
