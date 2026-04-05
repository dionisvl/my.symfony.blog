package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"api-go/internal/repository"
	"api-go/internal/repository/testhelper"
)

func TestPostRepository_FindPublishedPaginated(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "author@test.com")
	insertPost(t, tdb, "Published Post", "published-post", false, userID)
	insertPost(t, tdb, "Draft Post", "draft-post", true, userID)

	repo := repository.NewPostRepository(tdb.DB)

	posts, total, err := repo.FindPublishedPaginated(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Published Post", posts[0].Title)
	assert.Equal(t, "/storage/blog_images/no-image.png", posts[0].ImageURL)
}

func TestPostRepository_FindPublishedBySlug(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "author2@test.com")
	insertPost(t, tdb, "Hello World", "hello-world", false, userID)

	repo := repository.NewPostRepository(tdb.DB)

	post, err := repo.FindPublishedBySlug(ctx, "hello-world")
	require.NoError(t, err)
	assert.Equal(t, "Hello World", post.Title)
	assert.Equal(t, "hello-world", post.Slug)
}

func TestPostRepository_FindPublishedBySlug_NotFound(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	repo := repository.NewPostRepository(tdb.DB)

	_, err := repo.FindPublishedBySlug(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestPostRepository_FindFeatured(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "author3@test.com")
	insertFeaturedPost(t, tdb, "Featured", "featured-slug", userID)
	insertPost(t, tdb, "Regular", "regular-slug", false, userID)

	repo := repository.NewPostRepository(tdb.DB)

	posts, err := repo.FindFeatured(ctx, 3)
	require.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Featured", posts[0].Title)
}

func TestPostRepository_SearchPublished(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "author4@test.com")
	insertPost(t, tdb, "Go programming tips", "go-tips", false, userID)
	insertPost(t, tdb, "Python basics", "python-basics", false, userID)
	insertPost(t, tdb, "Go draft", "go-draft", true, userID)

	repo := repository.NewPostRepository(tdb.DB)

	posts, err := repo.SearchPublished(ctx, "go", 20)
	require.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Go programming tips", posts[0].Title)
}

func TestPostRepository_IncrementViews(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "author5@test.com")
	postID := insertPost(t, tdb, "View me", "view-me", false, userID)

	repo := repository.NewPostRepository(tdb.DB)

	err := repo.IncrementViews(ctx, postID)
	require.NoError(t, err)

	post, err := repo.FindPublishedBySlug(ctx, "view-me")
	require.NoError(t, err)
	assert.Equal(t, 1, post.ViewsCount)
}

func TestPostRepository_FindPublishedByCategorySlugPaginated(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "author6@test.com")
	catID := insertCategory(t, tdb, "Go", "go-category")
	insertPostWithCategory(t, tdb, "Go post", "go-post", false, userID, catID)
	insertPost(t, tdb, "Other post", "other-post", false, userID)

	repo := repository.NewPostRepository(tdb.DB)

	posts, total, err := repo.FindPublishedByCategorySlugPaginated(ctx, "go-category", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, posts, 1)
	assert.Equal(t, "Go post", posts[0].Title)
}

// --- helpers ---

func insertUser(t *testing.T, tdb *testhelper.TestDB, email string) int {
	t.Helper()
	var id int
	err := tdb.SQLDB().QueryRowContext(context.Background(),
		`INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		"Test Author", email, "hash", time.Now(), time.Now(),
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func insertPost(t *testing.T, tdb *testhelper.TestDB, title, slug string, draft bool, userID int) int {
	t.Helper()
	var id int
	err := tdb.SQLDB().QueryRowContext(context.Background(),
		`INSERT INTO posts (title, slug, status, is_featured, views_count, user_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		title, slug, draft, false, 0, userID, time.Now(), time.Now(),
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func insertFeaturedPost(t *testing.T, tdb *testhelper.TestDB, title, slug string, userID int) int {
	t.Helper()
	var id int
	err := tdb.SQLDB().QueryRowContext(context.Background(),
		`INSERT INTO posts (title, slug, status, is_featured, views_count, user_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		title, slug, false, true, 0, userID, time.Now(), time.Now(),
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func insertCategory(t *testing.T, tdb *testhelper.TestDB, title, slug string) int {
	t.Helper()
	var id int
	err := tdb.SQLDB().QueryRowContext(context.Background(),
		`INSERT INTO categories (title, slug, created_at, updated_at) VALUES ($1,$2,$3,$4) RETURNING id`,
		title, slug, time.Now(), time.Now(),
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func insertPostWithCategory(t *testing.T, tdb *testhelper.TestDB, title, slug string, draft bool, userID, catID int) int {
	t.Helper()
	var id int
	err := tdb.SQLDB().QueryRowContext(context.Background(),
		`INSERT INTO posts (title, slug, status, is_featured, views_count, user_id, category_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		title, slug, draft, false, 0, userID, catID, time.Now(), time.Now(),
	).Scan(&id)
	require.NoError(t, err)
	return id
}
