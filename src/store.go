package src

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Store struct {
	directoryPath string
}

func NewStore(directoryPath string) *Store {
	return &Store{
		directoryPath: directoryPath,
	}
}

func (store *Store) write(originFile io.Reader, fileName string) error {
	filePath := filepath.Join(store.directoryPath, fileName)
	destinationFile, err := os.Create(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("cant create originFile in %v, %v", store.directoryPath, err))
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, originFile); err != nil {
		return errors.New(fmt.Sprintf("can not copy origin to destination: %v", err))
	}
	return nil
}

func (store *Store) FileInfo(fileId string) (os.FileInfo, string, error) {
	fileNotFound := errors.New("file not found")
	filePath := filepath.Join(store.directoryPath, fileId)
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, "", fileNotFound
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, "", fileNotFound
	}
	return stat, filePath, nil
}
