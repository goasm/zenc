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
	writer1 := NewCryptoWriter(pass, header.IV[:], ofile)
	writer2 := NewChunkWriter(writer1)
	_, err = io.Copy(writer2, ifile)
	if err != nil {
		return err
	}
	err = writer2.Flush()
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
	reader1 := NewCryptoReader(pass, header.IV[:], ifile)
	reader2 := NewChunkReader(reader1)
	_, err = io.Copy(ofile, reader2)
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
