package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup

	writeToMap := func(offset, limit, index int64, filePath string, result [][]byte, wg *sync.WaitGroup, mu *sync.Mutex) {
		defer wg.Done()
		file, err := os.Open(filePath)
		if err != nil {

		}
		_, err = file.Seek(offset, io.SeekStart)
		if err != nil {

		}
		buffer := make([]byte, limit)
		_, err = file.Read(buffer)
		if err != nil {

		}
		//fmt.Printf("index %v; %v\n", index, string(buffer))
		mu.Lock()
		result[index] = buffer
		mu.Unlock()
	}
	filePath := ".\\tmp\\test.txt"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	var offset, size, limit, index int64 = 0, fileInfo.Size(), 5, 0
	result := make([][]byte, size/limit+1)
	for offset < size {
		wg.Add(1)
		go writeToMap(offset, limit, index, filePath, result, &wg, &mu)
		offset += limit
		index++
	}
	wg.Wait()
	readed := ""
	for _, d := range result {
		readed += string(d)
	}
	fmt.Println(readed)
}
