package handler

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var allowedSEOHosts = map[string]bool{
	"web3main.pro": true,
	"phpqa.ru":     true,
}

const defaultSEOHost = "web3main.pro"

type SEOHandler struct {
	seoDir string
	logger *slog.Logger
}

func NewSEOHandler(seoDir string, logger *slog.Logger) *SEOHandler {
	return &SEOHandler{seoDir: seoDir, logger: logger}
}

func normalizeSEOHost(raw string) string {
	if raw == "" {
		return defaultSEOHost
	}

	host := strings.ToLower(raw)

	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	host = strings.TrimPrefix(host, "www.")

	if !allowedSEOHosts[host] {
		return defaultSEOHost
	}

	return host
}

func (h *SEOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := normalizeSEOHost(r.Host)

	base := strings.TrimPrefix(r.URL.Path, "/")
	ext := filepath.Ext(base)
	fileType := strings.TrimSuffix(base, ext)
	extNoDot := strings.TrimPrefix(ext, ".")

	filename := fileType + "-" + host + "." + extNoDot
	fullPath := filepath.Join(h.seoDir, filename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		h.logger.Error("seo file read failed", slog.String("path", fullPath), slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	contentType := "text/plain; charset=utf-8"
	if extNoDot == "xml" {
		contentType = "application/xml; charset=utf-8"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}