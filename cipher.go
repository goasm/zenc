package zenc

import (
	"crypto/aes"
	"crypto/cipher"
)

func newBlockCipher(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}

type encrypter struct {
	s cipher.Stream
}

// NewEncrypter creates a Stage for encryption
func NewEncrypter(key []byte, iv []byte) Stage {
	block := newBlockCipher(key)
	return &encrypter{cipher.NewCFBEncrypter(block, iv)}
}

func (c *encrypter) Process(data []byte) []byte {
	c.s.XORKeyStream(data, data)
	return data
}

type decrypter struct {
	s cipher.Stream
}

// NewDecrypter creates a Stage for decryption
func NewDecrypter(key []byte, iv []byte) Stage {
	block := newBlockCipher(key)
	return &decrypter{cipher.NewCFBDecrypter(block, iv)}
}

func (c *decrypter) Process(data []byte) []byte {
	c.s.XORKeyStream(data, data)
	return data
}
