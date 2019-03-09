package zenc

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

// FileHeader is used to identify encrypted file
type FileHeader struct {
	Magic    uint32
	IV       [16]byte
	Checksum uint16
}

// NewFileHeader creates a new FileHeader
func NewFileHeader() *FileHeader {
	header := &FileHeader{
		Magic:    0x5A454E43, // ZENC
		Checksum: 0,
	}
	_, err := rand.Read(header.IV[:])
	if err != nil {
		panic(err)
	}
	return header
}

// ReadDataFrom reads the data of FileHeader from r
func (h *FileHeader) ReadDataFrom(r io.Reader) {
	err := binary.Read(r, binary.BigEndian, h)
	if err != nil {
		panic(err)
	}
}

// WriteDataTo writes the data of FileHeader to w
func (h *FileHeader) WriteDataTo(w io.Writer) {
	err := binary.Write(w, binary.BigEndian, h)
	if err != nil {
		panic(err)
	}
}
