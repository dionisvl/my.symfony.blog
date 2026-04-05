package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"api-go/internal/repository"
	"api-go/internal/repository/testhelper"
)

func TestSubscriptionRepository_Create(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "subscriptions")

	ctx := context.Background()
	repo := repository.NewSubscriptionRepository(tdb.DB)

	sub, err := repo.Create(ctx, "user@example.com", "test-token-uuid")
	require.NoError(t, err)
	assert.NotZero(t, sub.ID)
	assert.Equal(t, "user@example.com", sub.Email)
	assert.False(t, sub.CreatedAt.IsZero())
}

func TestSubscriptionRepository_EmailExists(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "subscriptions")

	ctx := context.Background()
	repo := repository.NewSubscriptionRepository(tdb.DB)

	_, err := repo.Create(ctx, "existing@example.com", "some-token")
	require.NoError(t, err)

	exists, err := repo.EmailExists(ctx, "existing@example.com")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.EmailExists(ctx, "unknown@example.com")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestSubscriptionRepository_Create_DuplicateEmail(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "subscriptions")

	ctx := context.Background()
	repo := repository.NewSubscriptionRepository(tdb.DB)

	_, err := repo.Create(ctx, "dup@example.com", "token-1")
	require.NoError(t, err)

	_, err = repo.Create(ctx, "dup@example.com", "token-2")
	assert.Error(t, err, "duplicate email should fail at DB level")
}
