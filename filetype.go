package zenc

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

// FileHeader is used to identify encrypted file
type FileHeader struct {
	Magic   uint32
	Version uint16
	IV      [16]byte
}

// NewFileHeader creates a new FileHeader
func NewFileHeader() *FileHeader {
	header := &FileHeader{
		Magic:   0x5A454E43, // ZENC
		Version: 0x0100,     // 1.0
	}
	_, err := rand.Read(header.IV[:])
	if err != nil {
		panic(err)
	}
	return header
}

// ReadFrom reads the data of FileHeader from r
func (h *FileHeader) ReadFrom(r io.Reader) (n int64, err error) {
	n = int64(binary.Size(h))
	err = binary.Read(r, binary.BigEndian, h)
	if err != nil {
		n = 0
	}
	return n, err
}

// WriteTo writes the data of FileHeader to w
func (h *FileHeader) WriteTo(w io.Writer) (n int64, err error) {
	n = int64(binary.Size(h))
	err = binary.Write(w, binary.BigEndian, h)
	if err != nil {
		n = 0
	}
	return n, err
}

// FileFooter is used to store extra file information
type FileFooter struct {
	Checksum uint32
}

// NewFileFooter creates a new FileFooter
func NewFileFooter() *FileFooter {
	return &FileFooter{0}
}
