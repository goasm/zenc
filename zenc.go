package zenc

import (
	"bytes"
	"io"
)

func verifyFileType(header FileHeader) error {
	if !bytes.Equal(header.Magic[:], Magic[:]) {
		return ErrWrongFileType
	}
	return nil
}

func verifyChecksum(expected, actual uint32) error {
	if expected != actual {
		return ErrWrongChecksum
	}
	return nil
}

// EncryptFile encrypts file using the given password
func EncryptFile(src io.Reader, dst io.Writer, pass string) (err error) {
	header := NewFileHeader()
	footer := NewFileFooter()
	pipeline := NewPipeline()
	_, err = header.WriteTo(dst)
	if err != nil {
		return
	}
	checksumStage := NewChecksumStage()
	pipeline.AddStage(checksumStage)
	pipeline.AddStage(NewCryptoStage(pass, header.IV[:]))
	pipeline.AddStage(NewChunkStage())
	pipeline.AddStage(NewDestStage(dst))
	_, err = io.Copy(pipeline, src)
	if err != nil {
		return
	}
	err = pipeline.Close()
	if err != nil {
		return
	}
	footer.Checksum = checksumStage.Sum
	_, err = footer.WriteTo(dst)
	return
}

// DecryptFile decrypts file using the given password
func DecryptFile(src io.Reader, dst io.Writer, pass string) (err error) {
	header := FileHeader{}
	footer := FileFooter{}
	pipeline := NewPipeline()
	_, err = header.ReadFrom(src)
	if err != nil {
		return
	}
	err = verifyFileType(header)
	if err != nil {
		return
	}
	checksumStage := NewChecksumStage()
	pipeline.AddStage(checksumStage)
	pipeline.AddStage(NewCryptoStage(pass, header.IV[:]))
	pipeline.AddStage(NewDechunkStage())
	pipeline.AddStage(NewSourceStage(src))
	_, err = io.Copy(dst, pipeline)
	if err != nil {
		return
	}
	err = pipeline.Close()
	if err != nil {
		return
	}
	_, err = footer.ReadFrom(src)
	if err != nil {
		return
	}
	err = verifyChecksum(footer.Checksum, checksumStage.Sum)
	return
}
