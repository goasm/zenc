package zenc

import (
	"io"
	"os"
)

func verifyChecksum(expected, actual uint32) error {
	if expected != actual {
		return ErrWrongChecksum
	}
	return nil
}

// EncryptFile encrypts file using the given password
func EncryptFile(ifile, ofile *os.File, pass string) error {
	header := NewFileHeader()
	_, err := header.WriteTo(ofile)
	if err != nil {
		return err
	}
	var writer io.Writer = ofile
	writer = NewCryptoWriter(pass, header.IV[:], writer)
	writer = NewChunkWriter(writer)
	_, err = io.Copy(writer, ifile)
	if err != nil {
		return err
	}
	footer := NewFileFooter()
	footer.Checksum = 0
	_, err = footer.WriteTo(ofile)
	if err != nil {
		return err
	}
	return nil
}

// DecryptFile decrypts file using the given password
func DecryptFile(ifile, ofile *os.File, pass string) error {
	header := FileHeader{}
	_, err := header.ReadFrom(ifile)
	if err != nil {
		return err
	}
	var reader io.Reader = ifile
	reader = NewChunkReader(reader)
	reader = NewCryptoReader(pass, header.IV[:], reader)
	_, err = io.Copy(ofile, reader)
	if err != nil {
		return err
	}
	footer := FileFooter{}
	_, err = footer.ReadFrom(ifile)
	if err != nil {
		return err
	}
	err = verifyChecksum(footer.Checksum, 0)
	if err != nil {
		return err
	}
	return nil
}
