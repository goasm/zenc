package zenc_test

import (
	"bytes"
	"testing"

	"github.com/goasm/zenc"
)

var inputData []byte

func init() {
	inputData = make([]byte, 5000)
	fillFib(inputData)
}

func fillFib(buf []byte) {
	a := 0
	b := 1
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(b)
		c := a + b
		a = b
		b = c
	}
}

func TestZencEncryption(t *testing.T) {
	src := bytes.NewBuffer(inputData)
	dst := new(bytes.Buffer)
	err := zenc.EncryptFile(src, dst, "passwd123")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func TestZencDecryption(t *testing.T) {
	src := bytes.NewBuffer(inputData)
	tmp := new(bytes.Buffer)
	dst := new(bytes.Buffer)
	err := zenc.EncryptFile(src, tmp, "passwd123")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	err = zenc.DecryptFile(tmp, dst, "passwd123")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if !bytes.Equal(dst.Bytes(), inputData) {
		t.Fatal("Decryption failed, mismatched data")
	}
}
