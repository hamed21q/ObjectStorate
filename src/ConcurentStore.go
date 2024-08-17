package src

import (
	"io"
	"log"
	"os"
)

type FilePart struct {
	ID   int
	Part []byte
}

func fileSize(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

const limit int = 2

func ReadPart(index int, offset int64, filePath string, result chan FilePart) {
	file, _ := os.Open(filePath)
	file.Seek(offset, io.SeekStart)
	buffer := make([]byte, limit)
	file.Read(buffer)
	result <- FilePart{
		Part: buffer,
		ID:   index,
	}
}

func ReadFile(filePath string) []byte {
	size, _ := fileSize(filePath)
	parts := make(chan FilePart)
	chunkCount := int(size)/limit + 1
	for i := 0; i < chunkCount; i++ {
		go ReadPart(i, int64(i*limit+1), filePath, parts)
	}
	buffer := make([]byte, size)
	for i := 0; i < chunkCount; i++ {
		part := <-parts
		copy(buffer[part.ID*limit:], part.Part)
	}
	return buffer
}

func WriteFile(data []byte, filePath string) error {
	chunkCount := len(data)/limit + 1
	errors := make(chan error)
	for i := 0; i < chunkCount; i++ {
		go func(data []byte, offset int64, errors chan error) {
			file, _ := os.Open(filePath)
			_, err := file.WriteAt(data, offset)
			if err != nil {
				errors <- err
			} else {
				errors <- nil
			}
		}(data[i*limit:(i+1)*limit], int64(i*limit), errors)
	}

	for i := 0; i < chunkCount; i++ {
		err := <-errors
		if err != nil {
			return err
		}
	}
	return nil
}
