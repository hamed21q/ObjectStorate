package main

import "ObjectStorage/src"

func main() {
	store := src.NewStore("./tmp")
	api := src.NewApi(store)
	api.Start(":8080")
}
