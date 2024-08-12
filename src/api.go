package src

import (
	"ObjectStorage/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
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
	ID     string
	Status int
}

func (api *Api) upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 30)
	if err != nil {
		http.Error(w, "unable to parse multipart", http.StatusInternalServerError)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "unable to cast multipart", http.StatusInternalServerError)
	}
	uniqueID := utils.GetUniqueID(header.Filename) + header.Filename
	if err = api.store.write(file, uniqueID); err != nil {
		http.Error(w, "failed to persist file", http.StatusInternalServerError)
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(UploadResponse{
		ID:     uniqueID,
		Status: 200,
	})
}

func (api *Api) download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["id"]

	filePath := filepath.Join(api.store.directoryPath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not retrieve file information.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	http.ServeFile(w, r, filePath)
}
