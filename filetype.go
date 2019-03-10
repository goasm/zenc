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
	Size    uint32
	IV      [16]byte
}

// NewFileHeader creates a new FileHeader
func NewFileHeader() *FileHeader {
	rd := [16]byte{}
	_, err := rand.Read(rd[:])
	if err != nil {
		panic(err)
	}
	return &FileHeader{
		Magic:   0x5A454E43, // ZENC
		Version: 0x0100,     // 1.0
		Size:    0,
		IV:      rd,
	}
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

// ReadFrom reads the data of FileFooter from r
func (f *FileFooter) ReadFrom(r io.Reader) (n int64, err error) {
	n = int64(binary.Size(f))
	err = binary.Read(r, binary.BigEndian, f)
	if err != nil {
		n = 0
	}
	return n, err
}

// WriteTo writes the data of FileFooter to w
func (f *FileFooter) WriteTo(w io.Writer) (n int64, err error) {
	n = int64(binary.Size(f))
	err = binary.Write(w, binary.BigEndian, f)
	if err != nil {
		n = 0
	}
	return n, err
}
