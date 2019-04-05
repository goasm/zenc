package zenc

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

// NewCryptoReader creates a Reader that reads ciphertext into plaintext
func NewCryptoReader(key []byte, iv []byte, r io.Reader) io.Reader {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, iv)
	return &cipher.StreamReader{S: stream, R: r}
}

// NewCryptoWriter creates a Writer that writes plaintext to ciphertext
func NewCryptoWriter(key []byte, iv []byte, w io.Writer) io.Writer {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, iv)
	return &cipher.StreamWriter{S: stream, W: w}
}
