package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	util "hieupc05.github/backend-server/internal/utils"
	"hieupc05.github/backend-server/internal/utils/random"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)

	userId := random.RandomInt(100000, 999999)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userId, util.DepositorRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotZero(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userId, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(random.RandomInt(100000, 999999), util.DepositorRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotZero(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
