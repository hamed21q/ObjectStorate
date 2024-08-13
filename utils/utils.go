package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetUniqueID(fileName string) string {
	currentTimestamp := time.Now().Format(time.RFC3339)
	concatenatedString := currentTimestamp + fileName
	hash := sha256.New()
	hash.Write([]byte(concatenatedString))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}

func RemoveMultipartForm(r *http.Request) {
	if r.MultipartForm != nil {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			log.Printf("error on removing multipart file: %v \n", err)
		}
	}
}

func FileFormatFromUrl(url string) string {
	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]
	fileFormat := strings.Split(fileName, ".")
	return "." + fileFormat[len(fileFormat)-1]
}
