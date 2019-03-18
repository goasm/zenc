package zenc

import (
	"crypto/sha256"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func keygen(pass string, len int) []byte {
	salt := []byte("zenc-v1.0")
	return pbkdf2.Key([]byte(pass), salt, 4096, len, sha256.New)
}

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
	checksum := uint32(0)
	pipeline := NewPipeline()
	pipeline.AddStage(NewChecksum(&checksum))
	pipeline.AddStage(NewEncrypter(keygen(pass, 32), header.IV[:]))
	pipeline.Run(ifile, ofile)
	footer := NewFileFooter()
	footer.Checksum = checksum
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
	checksum := uint32(0)
	pipeline := NewPipeline()
	pipeline.AddStage(NewDecrypter(keygen(pass, 32), header.IV[:]))
	pipeline.AddStage(NewChecksum(&checksum))
	pipeline.Run(ifile, ofile)
	footer := FileFooter{}
	_, err = footer.ReadFrom(ifile)
	if err != nil {
		return err
	}
	err = verifyChecksum(footer.Checksum, checksum)
	if err != nil {
		return err
	}
	return nil
}
