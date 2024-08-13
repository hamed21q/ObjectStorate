package src

import (
	"ObjectStorage/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
	"strings"
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
	api.router.HandleFunc("/upload", api.Upload)
	api.router.HandleFunc("/download/{id}", api.Download)
	err := http.ListenAndServe(address, api.router)
	if err != nil {
		log.Fatalf("can not run the server: %v", err)
	}
}

type UploadResponse struct {
	ID string
}

func (api *Api) Upload(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		api.MultipartUpload(w, r)
		return
	}

	if strings.HasPrefix(contentType, "application/json") {
		api.IndirectUploadWithUrl(w, r)
		return
	}

	http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
}

type UploadFileRequest struct {
	URL string `json:"file"`
}

func (api *Api) IndirectUploadWithUrl(w http.ResponseWriter, r *http.Request) {
	var req UploadFileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	file, err := api.store.DownloadFromUrl(&http.Client{}, req.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("can not Download url: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fileName := utils.FileFormatFromUrl(req.URL)
	fileID := utils.GetUniqueID(fileName) + fileName
	if err := api.store.write(file, fileID); err != nil {
		http.Error(w, fmt.Sprintf("failed to persist file: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(UploadResponse{
		ID: fileID,
	})
}

func (api *Api) MultipartUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "unable to cast multipart", http.StatusInternalServerError)
	}
	defer file.Close()
	defer utils.RemoveMultipartForm(r)
	uniqueID := utils.GetUniqueID(header.Filename) + header.Filename
	if err = api.store.write(file, uniqueID); err != nil {
		http.Error(w, "failed to persist file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UploadResponse{
		ID: uniqueID,
	})
}

func (api *Api) Download(w http.ResponseWriter, r *http.Request) {
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
