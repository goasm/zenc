package zenc

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	maxChunkSize = 4096
)

type ChunkStage struct {
	MiddleStage
	buffer bytes.Buffer
}

func NewChunkStage() *ChunkStage {
	return &ChunkStage{MiddleStage{}, bytes.Buffer{}}
}

func (cs *ChunkStage) Write(data []byte) (n int, err error) {
	n, err = cs.buffer.Write(data)
	if err != nil {
		return
	}
	prefix := [4]byte{}
	for cs.buffer.Len() >= maxChunkSize {
		binary.LittleEndian.PutUint32(prefix[:], uint32(maxChunkSize))
		_, err = cs.next.Write(prefix[:])
		if err != nil {
			break
		}
		_, err = cs.next.Write(cs.buffer.Next(maxChunkSize))
		if err != nil {
			break
		}
	}
	return
}

func (cs *ChunkStage) Flush() (err error) {
	prefix := [4]byte{}
	binary.LittleEndian.PutUint32(prefix[:], uint32(cs.buffer.Len()))
	_, err = cs.next.Write(prefix[:])
	if err != nil {
		return
	}
	_, err = cs.next.Write(cs.buffer.Bytes())
	if err != nil {
		return
	}
	return
}

type UnchunkStage struct {
	MiddleStage
	buffer bytes.Buffer
}

func NewUnchunkStage() *UnchunkStage {
	return &UnchunkStage{MiddleStage{}, bytes.Buffer{}}
}

func (us *UnchunkStage) Write(data []byte) (n int, err error) {
	n, err = us.buffer.Write(data)
	if err != nil {
		return
	}
	prefix := us.buffer.Bytes()[:4]
	chunkSize := int(binary.LittleEndian.Uint32(prefix))
	for us.buffer.Len() >= len(prefix)+chunkSize {
		us.buffer.Next(len(prefix))
		_, err = us.next.Write(us.buffer.Next(chunkSize))
		if err != nil {
			break
		}
		prefix = us.buffer.Bytes()[:4]
		chunkSize = int(binary.LittleEndian.Uint32(prefix))
	}
	return
}
