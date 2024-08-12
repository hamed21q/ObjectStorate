package src

import (
	"bytes"
	"encoding/json"
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
