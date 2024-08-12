package utils

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"mime/multipart"
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
