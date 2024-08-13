package src

import (
	"ObjectStorage/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
)

type Api struct {
	store  *Store
	router *mux.Router
}

func NewApi(store *Store) *Api {
	return &Api{
		store:  store,
		router: mux.NewRouter(),
	}
}

func (api *Api) Start(address string) {
	api.router.HandleFunc("/upload", api.upload)
	api.router.HandleFunc("/download/{id}", api.download)
	err := http.ListenAndServe(address, api.router)
	if err != nil {
		log.Fatalf("can not run the server: %v", err)
	}
}

type UploadResponse struct {
	ID string
}

func (api *Api) upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "unable to cast multipart", http.StatusInternalServerError)
	}
	defer file.Close()
	defer utils.RemoveMultipartForm(r)
	uniqueID := utils.GetUniqueID(header.Filename) + header.Filename
	if err = api.store.write(file, uniqueID); err != nil {
		http.Error(w, "failed to persist file", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UploadResponse{
		ID: uniqueID,
	})
}

func (api *Api) download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["id"]
	fileStat, filePath, err := api.store.FileInfo(fileID)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()))

	http.ServeFile(w, r, filePath)
}
