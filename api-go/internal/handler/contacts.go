package handler

import "net/http"

type ContactsHandler struct{}

func NewContactsHandler() *ContactsHandler {
	return &ContactsHandler{}
}

func (h *ContactsHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, map[string]string{
		"email":   "info@phpqa.ru",
		"website": "https://phpqa.ru",
	})
}