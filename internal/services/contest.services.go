package services

import (
	"context"
	"net/http"

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
	LiveContest()
	PlayContest()
	EndContest()
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

	contestCreateParams := db.CreateContestParams{
		UserID:      authPayload.UserID,
		SubjectID:   subjectID,
		NumQuestion: numQuestion,
		TimeExam:    timeExam,
		Questions:   MakeRandomQuestion(),
		State:       db.ContestStateIDLE,
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

func (c *ContestService) LiveContest() {}

func (c *ContestService) PlayContest() {}

func (c *ContestService) EndContest() {}

func (c *ContestService) SubmitAnswer() {}
