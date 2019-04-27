package zenc

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

// FileHeader is used to identify encrypted file
type FileHeader struct {
	Magic   [4]byte
	Version [2]byte
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
		Magic:   [4]byte{0x5A, 0x45, 0x4E, 0x43}, // ZENC
		Version: [2]byte{0x01, 0x00},             // 1.0
		IV:      rd,
	}
}

// ReadFrom reads the data of FileHeader from r
func (h *FileHeader) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.LittleEndian, h)
	if err == nil {
		n = int64(binary.Size(h))
	}
	return
}

// WriteTo writes the data of FileHeader to w
func (h *FileHeader) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.LittleEndian, h)
	if err == nil {
		n = int64(binary.Size(h))
	}
	return
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
	err = binary.Read(r, binary.LittleEndian, f)
	if err == nil {
		n = int64(binary.Size(f))
	}
	return
}

// WriteTo writes the data of FileFooter to w
func (f *FileFooter) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.LittleEndian, f)
	if err == nil {
		n = int64(binary.Size(f))
	}
	return
}
