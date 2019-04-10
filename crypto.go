package zenc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

func genKey(pass string, len int) []byte {
	salt := []byte("zenc-v1.0")
	return pbkdf2.Key([]byte(pass), salt, 4096, len, sha256.New)
}

// CryptoStage encrypts/decrypts input data using the cipher key-stream
type CryptoStage struct {
	stream cipher.Stream
	next   Stage
}

// NewCryptoStage creates a Stage for both encryption and decryption
func NewCryptoStage(key string, iv []byte) *CryptoStage {
	block, err := aes.NewCipher(genKey(key, 32))
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, iv)
	return &CryptoStage{stream, nil}
}

func (c *CryptoStage) SetNext(n Stage) {
	c.next = n
}

func (c *CryptoStage) Next() Stage {
	return c.next
}

func (c *CryptoStage) Write(data []byte) (int, error) {
	c.stream.XORKeyStream(data, data)
	return c.next.Write(data)
}
