package src

import (
	"fmt"
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

func ReadParts(chunkCount int, filePath string) <-chan FilePart {
	c := make(chan FilePart)
	for i := 0; i < chunkCount; i++ {
		go func(index int, offset int64) {
			file, _ := os.Open(filePath)
			file.Seek(offset, io.SeekStart)
			buffer := make([]byte, limit)
			file.Read(buffer)
			c <- FilePart{
				Part: buffer,
				ID:   index,
			}
		}(i, int64(i*limit+1))
	}
	return c
}

func ReadFile(filePath string) {
	size, _ := fileSize(filePath)
	chunkCount := int(size)/limit + 1
	parts := ReadParts(chunkCount, filePath)
	buffer := make([]byte, size)
	for i := 0; i < chunkCount; i++ {
		part := <-parts
		copy(buffer[part.ID*limit:], part.Part)
	}

	fmt.Println(string(buffer))
}

func main() {
	ReadFile(".\\tmp\\test.txt")
}
