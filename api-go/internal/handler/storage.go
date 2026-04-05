package handler

import "net/http"

func NewStorageHandler(storageDir string) http.Handler {
	return http.StripPrefix("/storage", http.FileServer(http.Dir(storageDir)))
}