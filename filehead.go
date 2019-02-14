package zenc

type FileHead struct {
	iv []byte
}

func NewFileHead() FileHead {
	return FileHead{}
}
