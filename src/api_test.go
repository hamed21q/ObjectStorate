package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func CreateRandomFile(t *testing.T) string {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	file, err := w.CreateFormFile("file", "testfile.txt")
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}

	fileContent := []byte("This is a test file.")
	_, err = file.Write(fileContent)
	assert.Nil(t, err)

	w.Close()

	req := httptest.NewRequest("POST", "/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	recorder := httptest.NewRecorder()

	api.upload(recorder, req)

	res := recorder.Result()

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.Nil(t, err)
	assert.NotNil(t, result["ID"].(string))
	_, err = os.Stat(filepath.Join(store.directoryPath, result["ID"].(string)))
	assert.False(t, os.IsExist(err))
	return result["ID"].(string)
}

func TestUploadHandler(t *testing.T) {
	CreateRandomFile(t)
}

func TestDownloadHandler(t *testing.T) {
	fileId := CreateRandomFile(t)
	req := httptest.NewRequest("GET", fmt.Sprintf("/download/%v", fileId), nil)
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/download/{id}", api.download)

	// Serve the request using the router
	router.ServeHTTP(recorder, req)
	api.download(recorder, req)

	res := recorder.Result()

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotNil(t, body)
	assert.Equal(t, string(body), "This is a test file.")
}

func TestDownloadHandlerFileNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", fmt.Sprintf("/download/%v", "some_not_exists_id"), nil)
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/download/{id}", api.download)

	router.ServeHTTP(recorder, req)
	api.download(recorder, req)

	res := recorder.Result()

	assert.Equal(t, res.StatusCode, http.StatusNotFound)
}
