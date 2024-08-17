package main

import (
	"os"
)

const limit int = 2

func WriteFile(data []byte, filePath string) error {
	chunkCount := len(data)/limit + 1
	errors := make(chan error)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	//file.WriteAt([]byte("hello hamed"), 5)
	for i := 0; i < chunkCount; i++ {
		go func(data []byte, offset int64, errors chan error) {
			file, _ := os.OpenFile(filePath, os.O_CREATE, 0660)
			defer file.Close()
			_, err := file.WriteAt(data, offset)
			if err != nil {
				errors <- err
			} else {
				errors <- nil
			}
		}(data[i*limit:], int64(i*limit), errors)
	}

	for i := 0; i < chunkCount; i++ {
		err := <-errors
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := WriteFile([]byte("hello how are you today?"), ".\\tmp\\test2.txt")
	if err != nil {
		panic(err)
	}
}
