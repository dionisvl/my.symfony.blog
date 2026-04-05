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

func TestPostLikeRepository_CreateAndDelete(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "like-author@test.com")
	postID := insertPost(t, tdb, "Like target", "like-target", false, userID)

	repo := repository.NewPostLikeRepository(tdb.DB)

	ip := "127.0.0.1"
	like, err := repo.Create(ctx, postID, &ip, nil)
	require.NoError(t, err)
	assert.NotZero(t, like.ID)
	assert.Equal(t, postID, like.PostID)
	assert.Equal(t, &ip, like.IP)

	err = repo.DeleteByPostAndTime(ctx, postID, like.CreatedAt)
	require.NoError(t, err)

	var count int
	err = tdb.SQLDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM posts_likes WHERE post_id = $1", postID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestPostLikeRepository_DeleteByPostAndTime_WrongTime(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "like-author2@test.com")
	postID := insertPost(t, tdb, "Like target 2", "like-target-2", false, userID)

	repo := repository.NewPostLikeRepository(tdb.DB)

	ip := "127.0.0.1"
	_, err := repo.Create(ctx, postID, &ip, nil)
	require.NoError(t, err)

	wrongTime := time.Now().Add(-48 * time.Hour)
	err = repo.DeleteByPostAndTime(ctx, postID, wrongTime)
	require.NoError(t, err)

	var count int
	err = tdb.SQLDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM posts_likes WHERE post_id = $1", postID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestPostLikeRepository_PostExists(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "like-author3@test.com")
	postID := insertPost(t, tdb, "Exists for like", "exists-for-like", false, userID)

	repo := repository.NewPostLikeRepository(tdb.DB)

	exists, err := repo.PostExists(ctx, postID)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.PostExists(ctx, 99999)
	require.NoError(t, err)
	assert.False(t, exists)
}
