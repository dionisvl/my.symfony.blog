package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"api-go/internal/repository"
)

type CommentHandler struct {
	comments repository.CommentRepository
	logger   *slog.Logger
}

func NewCommentHandler(comments repository.CommentRepository, logger *slog.Logger) *CommentHandler {
	return &CommentHandler{comments: comments, logger: logger}
}

type commentRequest struct {
	Message  string `json:"message"`
	PostID   int    `json:"post_id"`
	CountMe  int    `json:"countMe"`
	Honeypot string `json:"honeypot"`
}

func (h *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req commentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Message = strings.TrimSpace(req.Message)

	if req.Message == "" {
		RespondError(w, http.StatusUnprocessableEntity, "comment text is required")
		return
	}

	if req.CountMe < 3 {
		RespondError(w, http.StatusUnprocessableEntity, "anti-bot check failed")
		return
	}

	if req.Honeypot != "" {
		RespondError(w, http.StatusUnprocessableEntity, "anti-bot check failed")
		return
	}

	if req.PostID <= 0 {
		RespondError(w, http.StatusUnprocessableEntity, "invalid post_id")
		return
	}

	exists, err := h.comments.PostExists(r.Context(), req.PostID)
	if err != nil {
		h.logger.Error("post exists check failed", slog.String("error", err.Error()))
		RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !exists {
		RespondError(w, http.StatusNotFound, "post not found")
		return
	}

	comment, err := h.comments.Create(r.Context(), req.PostID, "anon", req.Message)
	if err != nil {
		h.logger.Error("comment create failed", slog.String("error", err.Error()))
		RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]any{
		"status":  "ok",
		"message": "your comment will be added soon",
		"id":      comment.ID,
	})
}
