package main

import (
	"errors"
	"io"
	"os"
)

var (
	VideoFileMaxChunkSize int = 1024 * 1024
)

type VideoFileByteCounter int64

type VideoFileReader struct {
	VideoFile       *os.File
	TotalBytes      VideoFileByteCounter
	ChunkSize       VideoFileByteCounter
	CurrentMinBytes VideoFileByteCounter
	/*
		Note: CurrentMaxBytes is calculated from
		length of byte array
	*/
}

func GetVideoFileReader(filename string, chunkSize VideoFileByteCounter) (vfr *VideoFileReader, err error) {
	vfr = new(VideoFileReader)

	if vfr.VideoFile, err = os.Open(filename); err != nil {
		return
	}

	/*
		Calculate total bytes
	*/

	/*
		Set ChunkSize, and set CurrentMinBytes
	*/
	vfr.CurrentMinBytes = 0
	vfr.ChunkSize = chunkSize

	return
}

type VFRCurrentChunk struct {
	Bytes    []byte
	MinByte  VideoFileByteCounter
	MaxByte  VideoFileByteCounter
	Length   int
	Finished bool
}

func (vfr *VideoFileReader) GetNextChunk() (res *VFRCurrentChunk, err error) {
	res = new(VFRCurrentChunk)
	res.MinByte = vfr.CurrentMinBytes
	res.Bytes = make([]byte, vfr.ChunkSize)
	numBytes, err := vfr.VideoFile.Read(res.Bytes)
	if errors.Is(err, io.EOF) {
		/*
			Previous iteration read all the way to EOF.
		*/
		res.Finished = true
		res.MaxByte = res.MinByte // convert int numBytes to int64
		return res, nil
	} else if err != nil {
		return
	} else {
		res.MaxByte = res.MinByte + VideoFileByteCounter(numBytes) // convert int numBytes to int64
	}

	res.Length = numBytes
	res.MaxByte = res.MinByte + VideoFileByteCounter(numBytes) // convert int numBytes to int64
	res.Finished = false
	vfr.CurrentMinBytes += VideoFileByteCounter(numBytes)
	return
}