package zenc

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

// NewCryptoReader creates a Reader that reads ciphertext into plaintext
func NewCryptoReader(key []byte, iv []byte) io.Reader {
}

// NewCryptoWriter creates a Writer that writes plaintext to ciphertext
func NewCryptoWriter(key []byte, iv []byte) io.Writer {
}
