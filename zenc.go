package zenc

import (
	"io"
	"os"
)

func verifyChecksum(expected, actual uint32) error {
	if expected != actual {
		return ErrWrongChecksum
	}
	return nil
}

// EncryptFile encrypts file using the given password
func EncryptFile(ifile, ofile *os.File, pass string) (err error) {
	header := NewFileHeader()
	footer := NewFileFooter()
	pipeline := NewPipeline()
	_, err = header.WriteTo(ofile)
	if err != nil {
		return
	}
	checksumStage := NewChecksumStage()
	pipeline.AddStage(checksumStage)
	pipeline.AddStage(NewCryptoStage(pass, header.IV[:]))
	pipeline.AddStage(NewChunkStage())
	pipeline.AddStage(NewDestStage(ofile))
	_, err = io.Copy(pipeline, ifile)
	if err != nil {
		return
	}
	footer.Checksum = checksumStage.Sum
	_, err = footer.WriteTo(ofile)
	return
}

// DecryptFile decrypts file using the given password
func DecryptFile(ifile, ofile *os.File, pass string) (err error) {
	header := FileHeader{}
	footer := FileFooter{}
	pipeline := NewPipeline()
	_, err = header.ReadFrom(ifile)
	if err != nil {
		return
	}
	checksumStage := NewChecksumStage()
	pipeline.AddStage(checksumStage)
	pipeline.AddStage(NewCryptoStage(pass, header.IV[:]))
	pipeline.AddStage(NewChunkStage())
	pipeline.AddStage(NewSourceStage(ifile))
	_, err = io.Copy(ofile, pipeline)
	if err != nil {
		return
	}
	_, err = footer.ReadFrom(ifile)
	if err != nil {
		return
	}
	err = verifyChecksum(footer.Checksum, checksumStage.Sum)
	return
}
