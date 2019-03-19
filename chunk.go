package zenc

type chunkEncoder struct {
}

func NewChunkEncoder() Stage {
	return &chunkEncoder{}
}

func (c *chunkEncoder) Process(data []byte) []byte {
	return data
}

type chunkDecoder struct {
}

func NewChunkDecoder() Stage {
	return &chunkDecoder{}
}

func (c *chunkDecoder) Process(data []byte) []byte {
	return data
}
