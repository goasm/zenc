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

// EncryptFile encrypts file using the given password
func EncryptFile(ifile, ofile *os.File, pass string) {
	iv := make([]byte, 16)
	pipeline := NewPipeline()
	pipeline.AddStage(NewEncrypter(keygen(pass, 32), iv))
	pipeline.Run(ifile, ofile)
}

// DecryptFile decrypts file using the given password
func DecryptFile(ifile, ofile *os.File, pass string) {
	iv := make([]byte, 16)
	pipeline := NewPipeline()
	pipeline.AddStage(NewDecrypter(keygen(pass, 32), iv))
	pipeline.Run(ifile, ofile)
}
