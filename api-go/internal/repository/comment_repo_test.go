package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"api-go/internal/repository"
	"api-go/internal/repository/testhelper"
)

func TestCommentRepository_Create(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "comment-author@test.com")
	postID := insertPost(t, tdb, "Comment target", "comment-target", false, userID)

	repo := repository.NewCommentRepository(tdb.DB)

	comment, err := repo.Create(ctx, postID, "anon", "Hello, world!")
	require.NoError(t, err)
	assert.NotZero(t, comment.ID)
	assert.Equal(t, "anon", comment.AuthorName)
	assert.Equal(t, "Hello, world!", comment.Text)
	assert.Equal(t, 0, comment.Status)
}

func TestCommentRepository_PostExists(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "post_tags", "posts_likes", "comments", "posts", "users", "categories")

	ctx := context.Background()
	userID := insertUser(t, tdb, "exists-author@test.com")
	postID := insertPost(t, tdb, "Exists post", "exists-post", false, userID)

	repo := repository.NewCommentRepository(tdb.DB)

	exists, err := repo.PostExists(ctx, postID)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.PostExists(ctx, 99999)
	require.NoError(t, err)
	assert.False(t, exists)
}
