package src

import (
	"bytes"
	io "io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	fileContent := "This is a test file."
	fileBytes := []byte(fileContent)
	file := io.NopCloser(bytes.NewReader(fileBytes))

	fileName := "testFile.txt"

	err := store.write(file, fileName)
	require.NoError(t, err)

	expectedFilePath := filepath.Join(store.directoryPath, fileName)
	assert.FileExists(t, expectedFilePath)

	writtenContent, err := os.ReadFile(expectedFilePath)
	require.NoError(t, err)
	require.Equal(t, fileContent, string(writtenContent))
}
