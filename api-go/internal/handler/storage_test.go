package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageHandler(t *testing.T) {
	tempDir := t.TempDir()
	subdir := filepath.Join(tempDir, "subdir")
	require.NoError(t, os.Mkdir(subdir, 0755))

	testFile := filepath.Join(tempDir, "test.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("test content"), 0644))

	subFile := filepath.Join(subdir, "subfile.txt")
	require.NoError(t, os.WriteFile(subFile, []byte("sub content"), 0644))

	handler := NewStorageHandler(tempDir)

	t.Run("GET /storage/test.txt returns 200 with cache headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/storage/test.txt", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test content", w.Body.String())
		assert.Equal(t, "public, max-age=86400", w.Header().Get("Cache-Control"))
	})

	t.Run("GET /storage/ returns 404, no listing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/storage/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NotContains(t, w.Body.String(), "test.txt")
		assert.NotContains(t, w.Body.String(), "subdir")
	})

	t.Run("GET /storage/subdir/ returns 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/storage/subdir/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("GET /storage/subdir/subfile.txt returns 200 with cache headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/storage/subdir/subfile.txt", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "sub content", w.Body.String())
		assert.Equal(t, "public, max-age=86400", w.Header().Get("Cache-Control"))
	})

	t.Run("GET /storage/nonexistent returns 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/storage/nonexistent.txt", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
