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
func EncryptFile(ifile, ofile *os.File, pass string) error {
	header := NewFileHeader()
	footer := NewFileFooter()
	pipeline := NewPipeline()
	checksumStage := NewChecksumStage()
	pipeline.AddStage(checksumStage)
	pipeline.AddStage(NewCryptoStage(pass, header.IV[:]))
	pipeline.AddStage(NewChunkStage())
	pipeline.AddStage(NewDestStage(ofile))
	_, err := header.WriteTo(ofile)
	if err != nil {
		return err
	}
	_, err = io.Copy(pipeline, ifile)
	if err != nil {
		return err
	}
	footer.Checksum = checksumStage.Sum
	_, err = footer.WriteTo(ofile)
	if err != nil {
		return err
	}
	return nil
}

// DecryptFile decrypts file using the given password
func DecryptFile(ifile, ofile *os.File, pass string) error {
	header := FileHeader{}
	footer := FileFooter{}
	pipeline := NewPipeline()
	checksumStage := NewChecksumStage()
	_, err := header.ReadFrom(ifile)
	if err != nil {
		return err
	}
	pipeline.AddStage(checksumStage)
	pipeline.AddStage(NewCryptoStage(pass, header.IV[:]))
	pipeline.AddStage(NewChunkStage())
	pipeline.AddStage(NewSourceStage(ifile))
	_, err = io.Copy(ofile, pipeline)
	if err != nil {
		return err
	}
	_, err = footer.ReadFrom(ifile)
	if err != nil {
		return err
	}
	err = verifyChecksum(footer.Checksum, checksumStage.Sum)
	if err != nil {
		return err
	}
	return nil
}
