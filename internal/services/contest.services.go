package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hieupc05.github/backend-server/internal/utils/sse"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/middlewares"
	"hieupc05.github/backend-server/internal/utils/random"
	"hieupc05.github/backend-server/internal/utils/token"
	"hieupc05.github/backend-server/response"
)

type IContestService interface {
	CreateContest(ctx *gin.Context, subjectID int64, numQuestion int32, timeExam int32)
	LiveContest(ctx *gin.Context, state int64)
	PlayContest(ctx *gin.Context, rid int64)
	EndContest(ctx *gin.Context, rid int64)
	SubmitAnswer()
}

type ContestService struct{}

func NewContestService() *ContestService {
	return &ContestService{}
}

func MakeRandomQuestion() string {
	return random.RandomString(10)
}

func (c *ContestService) CreateContest(ctx *gin.Context, subjectID int64, numQuestion int32, timeExam int32) {
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	arg := db.GetContestInTwoCaseParams{
		UserID:  authPayload.UserID,
		State:   db.ContestStateRUNNING,
		State_2: db.ContestStateIDLE,
	}
	contestExists, err := global.PgDb.GetContestInTwoCase(context.Background(), arg)
	if err != nil {
		fmt.Println(err)
		rsp := response.ErrorResponse(response.ErrSystem)
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	if contestExists {
		rsp := response.ErrorResponse(response.ErrContestAlreadyExists)
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	contestCreateParams := db.CreateContestParams{
		UserID:        authPayload.UserID,
		TimeStartExam: time.Now().Unix(),
		SubjectID:     subjectID,
		NumQuestion:   numQuestion,
		TimeExam:      timeExam,
		Questions:     MakeRandomQuestion(),
		State:         db.ContestStateIDLE,
	}

	contestId, err := global.PgDb.CreateContest(context.Background(), contestCreateParams)
	if err != nil {
		rsp := response.ErrorResponse(response.ErrSystem)
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	rsp := response.SuccessResponse(response.SuccessCode, gin.H{"id": contestId})
	ctx.JSON(http.StatusOK, rsp)
}
func (c *ContestService) LiveContest(ctx *gin.Context, id int64) {
	authPayload, ok := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
	if !ok {
		rsp := response.ErrorResponse(response.ErrUnauthorizedInvalidToken)
		ctx.JSON(http.StatusUnauthorized, rsp)
		return
	}

	arg := db.GetUserContestByIDParams{
		ID:     id,
		UserID: authPayload.UserID,
	}
	fmt.Println(arg)
	// Fetch the contest from the database
	contest, err := global.PgDb.GetUserContestByID(context.Background(), arg)
	fmt.Println(contest)
	if err != nil {
		statusCode, rsp := handleContestFetchError(err)
		ctx.JSON(statusCode, response.ErrorResponse(rsp))
		return
	}

	// Validate contest state
	if contest.State == db.ContestStateFINISHED {
		rsp := response.ErrorResponse(response.ErrContestFinished)
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	if contest.State != db.ContestStateIDLE {
		rsp := response.ErrorResponse(response.ErrContestRunning)
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	// Create a contest room
	global.RoomManage.MakeRoom(id)

	rsp := response.SuccessResponse(response.SuccessCode, gin.H{"id": id})
	ctx.JSON(http.StatusOK, rsp)
}

// Helper function to handle contest fetch errors
func handleContestFetchError(err error) (int, int) {
	if errors.Is(err, db.ErrRecordNotFound) {
		return http.StatusBadRequest, response.ErrInvalidContestNotFound
	}
	return http.StatusInternalServerError, response.ErrSystem
}

func (c *ContestService) PlayContest(ctx *gin.Context, rid int64) {
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	arg := db.GetUserContestByIDParams{
		ID:     rid,
		UserID: authPayload.UserID,
	}
	contest, err := global.PgDb.GetUserContestByID(ctx, arg)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rsp := response.ErrorResponse(response.ErrInvalidContestNotFound)
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		rsp := response.ErrorResponse(response.ErrSystem)
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	if contest.State != db.ContestStateIDLE {
		rsp := response.ErrorResponse(response.ErrContestRunning)
		if contest.State == db.ContestStateFINISHED {
			rsp = response.ErrorResponse(response.ErrContestFinished)
		}
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	updateArg := db.UpdateContestStateParams{
		ID:    contest.ID,
		State: db.ContestStateRUNNING,
	}
	result, err := global.PgDb.UpdateContestState(context.Background(), updateArg)

	if err != nil {
		rsp := response.ErrorResponse(response.ErrSystem)
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}
	global.RoomManage.BroadcastToRoom(result.ID, "Start Contest")

	ctx.JSON(http.StatusOK, gin.H{"id": result.ID})
}

func (c *ContestService) EndContest(ctx *gin.Context, rid int64) {
	// TODO: Update contest state to FINISHED

	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	arg := db.GetUserContestByIDParams{
		ID:     rid,
		UserID: authPayload.UserID,
	}
	contest, err := global.PgDb.GetUserContestByID(ctx, arg)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rsp := response.ErrorResponse(response.ErrInvalidContestNotFound)
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		rsp := response.ErrorResponse(response.ErrSystem)
		ctx.JSON(http.StatusInternalServerError, rsp)
		return
	}

	if contest.State == db.ContestStateFINISHED {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrContestFinished))
		return
	}

	if roomNotExists := global.RoomManage.IsRoomNotExist(rid); roomNotExists {
		rsp := response.ErrorResponse(response.ErrInvalidData)
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	// Broadcast to all clients
	global.RoomManage.BroadcastToRoom(rid, sse.EndContest)
	global.RoomManage.Rooms[rid].State <- true

	//global.RoomManage.RemoveRoom(rid)
	ctx.JSON(http.StatusOK, gin.H{"id": rid})
}

func (c *ContestService) SubmitAnswer() {}
