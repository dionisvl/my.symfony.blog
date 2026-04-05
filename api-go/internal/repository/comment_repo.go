package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"api-go/internal/model"
)

type CommentRepository interface {
	Create(ctx context.Context, postID int, authorName, text string) (*model.Comment, error)
	PostExists(ctx context.Context, postID int) (bool, error)
}

type commentRepository struct {
	db *bun.DB
}

func NewCommentRepository(db *bun.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, postID int, authorName, text string) (*model.Comment, error) {
	now := time.Now()
	comment := &model.Comment{
		PostID:     postID,
		AuthorName: authorName,
		Text:       text,
		Status:     0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err := r.db.NewInsert().Model(comment).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *commentRepository) PostExists(ctx context.Context, postID int) (bool, error) {
	exists, err := r.db.NewSelect().
		TableExpr("posts").
		Where("id = ?", postID).
		Exists(ctx)
	return exists, err
}
