package zenc

import (
	"crypto/sha256"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func keygen(pass string, len int) []byte {
	salt := []byte("zenc-v1.0")
	return pbkdf2.Key([]byte(pass), salt, 4096, len, sha256.New)
}

// EncryptFile encrypts file using the given password
func EncryptFile(ifile, ofile *os.File, pass string) error {
	info, err := ifile.Stat()
	if err != nil {
		return err
	}
	header := NewFileHeader()
	header.Size = uint32(info.Size())
	_, err = header.WriteTo(ofile)
	if err != nil {
		return err
	}
	pipeline := NewPipeline()
	pipeline.AddStage(NewEncrypter(keygen(pass, 32), header.IV[:]))
	pipeline.Run(ifile, ofile)
	footer := NewFileFooter()
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
	pipeline := NewPipeline()
	pipeline.AddStage(NewDecrypter(keygen(pass, 32), header.IV[:]))
	pipeline.Run(io.LimitReader(ifile, int64(header.Size)), ofile)
	footer := FileFooter{}
	_, err = footer.ReadFrom(ifile)
	if err != nil {
		return err
	}
	return nil
}
