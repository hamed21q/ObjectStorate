package src

import (
	"bytes"
	io "io"
	"net/http"
	"net/http/httptest"
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

func TestDownloadURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "mock file content")
	}))
	defer server.Close()

	client := server.Client()

	file, err := store.DownloadFromUrl(client, server.URL)
	defer file.Close()
	assert.Nil(t, err)
	res, err := io.ReadAll(file)
	assert.Nil(t, err)
	assert.Equal(t, string(res), "mock file content")
}

func TestDownloadURLServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failed to persist file", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := server.Client()

	file, err := store.DownloadFromUrl(client, server.URL)
	assert.Nil(t, file)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "bad status: 500 Internal Server Error")
}
