package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"api-go/internal/model"
)

type PostLikeRepository interface {
	Create(ctx context.Context, postID int, ip *string, deviceMemory *int) (*model.PostLike, error)
	DeleteByIDAndPost(ctx context.Context, likeID, postID int) error
	PostExists(ctx context.Context, postID int) (bool, error)
}

type postLikeRepository struct {
	db *bun.DB
}

func NewPostLikeRepository(db *bun.DB) PostLikeRepository {
	return &postLikeRepository{db: db}
}

func (r *postLikeRepository) Create(ctx context.Context, postID int, ip *string, deviceMemory *int) (*model.PostLike, error) {
	now := time.Now()
	like := &model.PostLike{
		PostID:       postID,
		IP:           ip,
		DeviceMemory: deviceMemory,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	_, err := r.db.NewInsert().Model(like).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return like, nil
}

func (r *postLikeRepository) DeleteByIDAndPost(ctx context.Context, likeID, postID int) error {
	_, err := r.db.NewDelete().
		TableExpr("posts_likes").
		Where("id = ?", likeID).
		Where("post_id = ?", postID).
		Exec(ctx)
	return err
}

func (r *postLikeRepository) PostExists(ctx context.Context, postID int) (bool, error) {
	exists, err := r.db.NewSelect().
		TableExpr("posts").
		Where("id = ?", postID).
		Exists(ctx)
	return exists, err
}
