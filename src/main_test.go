package src

import (
	"log"
	"os"
	"testing"
)

var store *Store
var api *Api

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "store_test")
	if err != nil {
		log.Fatal(err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Fatal(err)
		}
	}(tmpDir)
	store = NewStore(tmpDir)
	api = NewApi(store)
	os.Exit(m.Run())
}
