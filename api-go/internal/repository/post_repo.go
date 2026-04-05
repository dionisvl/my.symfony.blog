package repository

import (
	"context"
	"strings"

	"github.com/uptrace/bun"

	"api-go/internal/model"
)

type PostRepository interface {
	FindPublishedPaginated(ctx context.Context, page, perPage int) ([]model.Post, int, error)
	FindPublishedBySlug(ctx context.Context, slug string) (*model.Post, error)
	FindFeatured(ctx context.Context, limit int) ([]model.Post, error)
	FindRecentPublished(ctx context.Context, limit int) ([]model.Post, error)
	FindPublishedByCategorySlugPaginated(ctx context.Context, slug string, page, perPage int) ([]model.Post, int, error)
	FindPublishedByTagSlugPaginated(ctx context.Context, slug string, page, perPage int) ([]model.Post, int, error)
	SearchPublished(ctx context.Context, query string, limit int) ([]model.Post, error)
	IncrementViews(ctx context.Context, postID int) error
}

type postRepository struct {
	db *bun.DB
}

func NewPostRepository(db *bun.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) FindPublishedPaginated(ctx context.Context, page, perPage int) ([]model.Post, int, error) {
	var posts []model.Post
	count, err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Where("p.status = ?", false).
		OrderExpr("p.created_at DESC").
		Limit(perPage).
		Offset((page - 1) * perPage).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	resolveImages(posts)
	return posts, count, nil
}

func (r *postRepository) FindPublishedBySlug(ctx context.Context, slug string) (*model.Post, error) {
	var post model.Post
	err := r.db.NewSelect().
		Model(&post).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Where("p.slug = ?", slug).
		Where("p.status = ?", false).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	comments, err := r.findApprovedComments(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments = comments

	likesCount, err := r.countLikes(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	post.LikesCount = likesCount

	post.ResolveImageURL()
	return &post, nil
}

func (r *postRepository) FindFeatured(ctx context.Context, limit int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Where("p.is_featured = ?", true).
		Where("p.status = ?", false).
		Limit(limit).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	resolveImages(posts)
	return posts, nil
}

func (r *postRepository) FindRecentPublished(ctx context.Context, limit int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Where("p.status = ?", false).
		OrderExpr("p.created_at DESC").
		Limit(limit).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	resolveImages(posts)
	return posts, nil
}

func (r *postRepository) FindPublishedByCategorySlugPaginated(ctx context.Context, slug string, page, perPage int) ([]model.Post, int, error) {
	var posts []model.Post
	count, err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Join("JOIN categories AS cat ON cat.id = p.category_id").
		Where("cat.slug = ?", slug).
		Where("p.status = ?", false).
		OrderExpr("p.created_at DESC").
		Limit(perPage).
		Offset((page - 1) * perPage).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	resolveImages(posts)
	return posts, count, nil
}

func (r *postRepository) FindPublishedByTagSlugPaginated(ctx context.Context, slug string, page, perPage int) ([]model.Post, int, error) {
	var posts []model.Post
	count, err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Relation("Tags").
		Join("JOIN post_tags AS pt ON pt.post_id = p.id").
		Join("JOIN tags AS tg ON tg.id = pt.tag_id").
		Where("tg.slug = ?", slug).
		Where("p.status = ?", false).
		OrderExpr("p.created_at DESC").
		Limit(perPage).
		Offset((page - 1) * perPage).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	resolveImages(posts)
	return posts, count, nil
}

func (r *postRepository) SearchPublished(ctx context.Context, query string, limit int) ([]model.Post, error) {
	escaped := escapeLike(query)
	pattern := "%" + escaped + "%"

	var posts []model.Post
	err := r.db.NewSelect().
		Model(&posts).
		Relation("Author").
		Relation("Category").
		Where("p.status = ?", false).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.
				Where("p.title ILIKE ?", pattern).
				WhereOr("p.description ILIKE ?", pattern).
				WhereOr("p.content ILIKE ?", pattern)
		}).
		OrderExpr("p.created_at DESC").
		Limit(limit).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	resolveImages(posts)
	return posts, nil
}

func (r *postRepository) IncrementViews(ctx context.Context, postID int) error {
	_, err := r.db.NewUpdate().
		TableExpr("posts").
		Set("views_count = views_count + 1").
		Where("id = ?", postID).
		Exec(ctx)
	return err
}

func (r *postRepository) findApprovedComments(ctx context.Context, postID int) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.NewSelect().
		Model(&comments).
		Where("post_id = ?", postID).
		Where("status = ?", 1).
		OrderExpr("created_at ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *postRepository) countLikes(ctx context.Context, postID int) (int, error) {
	count, err := r.db.NewSelect().
		TableExpr("posts_likes").
		Where("post_id = ?", postID).
		Count(ctx)
	return count, err
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func resolveImages(posts []model.Post) {
	for i := range posts {
		posts[i].ResolveImageURL()
	}
}