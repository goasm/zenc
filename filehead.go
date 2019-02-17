package zenc

// FileHead is used to identify encrypted file
type FileHead struct {
	magic    uint32
	iv       [16]byte
	checksum uint16
}
