package zenc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func genKey(pass string, len int) []byte {
	salt := []byte("zenc-v1.0")
	return pbkdf2.Key([]byte(pass), salt, 4096, len, sha256.New)
}

// NewCryptoReader creates a Reader that reads ciphertext into plaintext
func NewCryptoReader(key string, iv []byte, r io.Reader) io.Reader {
	block, err := aes.NewCipher(genKey(key, 32))
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, iv)
	return &cipher.StreamReader{S: stream, R: r}
}

// NewCryptoWriter creates a Writer that writes plaintext to ciphertext
func NewCryptoWriter(key string, iv []byte, w io.Writer) io.Writer {
	block, err := aes.NewCipher(genKey(key, 32))
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, iv)
	return &cipher.StreamWriter{S: stream, W: w}
}
