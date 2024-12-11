package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/internal/utils/random"
)

func CreateUser(t *testing.T) db.User {
	arg := db.CreateUserParams{
		Email:          random.RandomEmail(),
		HashedPassword: random.RandomString(8),
	}

	user, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateUser(t)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := CreateUser(t)

	user2, err := testStore.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}

func TestChangePassword(t *testing.T) {
	user1 := CreateUser(t)

	arg := db.ChangePasswordParams{
		Email:          user1.Email,
		HashedPassword: random.RandomString(8),
	}

	user2, err := testStore.ChangePassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
}
