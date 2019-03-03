package zenc

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

// FileHead is used to identify encrypted file
type FileHead struct {
	Magic    uint32
	IV       [16]byte
	Checksum uint16
}

// NewFileHead creates a new FileHead
func NewFileHead() FileHead {
	head := FileHead{
		Magic:    0x5A454E43, // ZENC
		Checksum: 0,
	}
	_, err := rand.Read(head.IV[:])
	if err != nil {
		panic(err)
	}
	return head
}

// ReadDataFrom reads the data of FileHead from r
func (h FileHead) ReadDataFrom(r io.Reader) {
	err := binary.Read(r, binary.BigEndian, &h)
	if err != nil {
		panic(err)
	}
}

// WriteDataTo writes the data of FileHead to w
func (h FileHead) WriteDataTo(w io.Writer) {
	err := binary.Write(w, binary.BigEndian, &h)
	if err != nil {
		panic(err)
	}
}
