package zenc

import (
	"os"
)

// EncryptFile encrypts file using the given password
func EncryptFile(ifile, ofile *os.File, pass string) {
	pipeline := NewPipeline()
	pipeline.AddStage(NewEncrypter())
	pipeline.Run(ifile, ofile)
}

// DecryptFile decrypts file using the given password
func DecryptFile(ifile, ofile *os.File, pass string) {
	pipeline := NewPipeline()
	pipeline.AddStage(NewDecrypter())
	pipeline.Run(ifile, ofile)
}
