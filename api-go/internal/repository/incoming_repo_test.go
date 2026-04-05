package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"api-go/internal/repository"
	"api-go/internal/repository/testhelper"
)

func TestIncomingRepository_Create(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "incomings")

	ctx := context.Background()
	repo := repository.NewIncomingRepository(tdb.DB)

	incoming, err := repo.Create(ctx, "Alice", "alice@example.com", "+79001234567", "Hello!")
	require.NoError(t, err)
	assert.NotZero(t, incoming.ID)
	assert.Equal(t, "Alice", incoming.Name)
	assert.Equal(t, "alice@example.com", incoming.Email)
	assert.Equal(t, "+79001234567", incoming.Phone)
	assert.Equal(t, "Hello!", incoming.Message)
	assert.Equal(t, 0, incoming.Status)
	assert.False(t, incoming.CreatedAt.IsZero())
}

func TestIncomingRepository_Create_EmptyOptionalFields(t *testing.T) {
	tdb := testhelper.NewTestDB(t)
	tdb.Truncate(t, "incomings")

	ctx := context.Background()
	repo := repository.NewIncomingRepository(tdb.DB)

	incoming, err := repo.Create(ctx, "Bob", "", "", "Just a message")
	require.NoError(t, err)
	assert.NotZero(t, incoming.ID)
	assert.Equal(t, "Bob", incoming.Name)
	assert.Equal(t, "", incoming.Email)
	assert.Equal(t, "", incoming.Phone)
	assert.Equal(t, 0, incoming.Status)
}
