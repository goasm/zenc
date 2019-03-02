package zenc

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

// FileHead is used to identify encrypted file
type FileHead struct {
	magic    uint32
	iv       [16]byte
	checksum uint16
}

// NewFileHead creates a new FileHead
func NewFileHead() FileHead {
	head := FileHead{
		magic:    0x5A454E43, // ZENC
		checksum: 0,
	}
	_, err := rand.Read(head.iv[:])
	if err != nil {
		panic(err)
	}
	return head
}

// ReadFrom reads the data of FileHead from r
func (h FileHead) ReadFrom(r io.Reader) {
	err := binary.Read(r, binary.BigEndian, &h)
	if err != nil {
		panic(err)
	}
}

// WriteTo writes the data of FileHead to w
func (h FileHead) WriteTo(w io.Writer) {
	err := binary.Write(w, binary.BigEndian, &h)
	if err != nil {
		panic(err)
	}
}
