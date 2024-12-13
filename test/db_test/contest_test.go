package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/internal/utils/random"
)

func CreateContest(t *testing.T) int64 {
	user := CreateUser(t)

	arg := db.CreateContestParams{
		UserID:        user.ID,
		SubjectID:     random.RandomInt(1, 10),
		NumQuestion:   int32(random.RandomInt(1, 10)),
		TimeExam:      int32(random.RandomInt(1, 10)),
		TimeStartExam: time.Now().Unix(),
		State:         db.ContestStateIDLE,
		Questions:     random.RandomString(10),
	}

	contestid, err := testStore.CreateContest(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, contestid)

	return contestid
}

func TestCreateContest(t *testing.T) {
	CreateContest(t)
}

func TestGetContest(t *testing.T) {
	contest := CreateContest(t)

	contestGet, err := testStore.GetContest(context.Background(), contest)
	require.NoError(t, err)
	require.NotEmpty(t, contestGet)

	require.Equal(t, contest, contestGet.ID)
	require.Equal(t, db.ContestStateIDLE, contestGet.State)
}

func TestUpdateContest(t *testing.T) {
	contest := CreateContest(t)
	timeExamNew := int32(random.RandomInt(1, 10))
	timeStartExamNew := time.Now().Unix()
	numQ := int32(random.RandomInt(1, 10))
	newQuestion := random.RandomString(15)
	stateNew := db.ContestStateRUNNING

	arg := db.UpdateContestParams{
		ID:            contest,
		NumQuestion:   numQ,
		TimeExam:      timeExamNew,
		TimeStartExam: timeStartExamNew,
		State:         stateNew,
		Questions:     newQuestion,
	}

	contestUpdate, err := testStore.UpdateContest(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contestUpdate)

	require.Equal(t, contest, contestUpdate.ID)
	require.Equal(t, numQ, contestUpdate.NumQuestion)
	require.Equal(t, timeExamNew, contestUpdate.TimeExam)
	require.Equal(t, stateNew, contestUpdate.State)
	require.Equal(t, newQuestion, contestUpdate.Questions)
}

func TestGetListContestByState(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateContest(t)
	}

	arg := db.ContestStateIDLE

	contestList, err := testStore.GetContestByState(context.Background(), arg)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(contestList), 10)

	for _, contest := range contestList {
		require.Equal(t, arg, contest.State)
	}
}
