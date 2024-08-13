package utils

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

func SafeClose(file multipart.File) {
	err := file.Close()
	if err != nil {
		log.Fatalf("cant close the file")
	}
}

func GetUniqueID(fileName string) string {
	currentTimestamp := time.Now().Format(time.RFC3339)
	concatenatedString := currentTimestamp + fileName
	hash := md5.New()
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
