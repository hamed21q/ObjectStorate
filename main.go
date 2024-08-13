package main

import (
	"ObjectStorage/src"
	"ObjectStorage/utils"
	"log"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	store := src.NewStore(config.FileStoragePath)
	api := src.NewApi(store)
	api.Start(config.HTTPServerAddress)
}
