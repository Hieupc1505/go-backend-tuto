package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hieupc05.github/backend-server/internal/services"
	"hieupc05.github/backend-server/response"
)

type ContestController struct {
	contestService services.IContestService
}

type CreateContestRequest struct {
	NumQuestion int32  `json:"num_question" binding:"required"`
	SubjectID   int64  `json:"subject_id" binding:"required"`
	SubjectName string `json:"subject_name"`
	TimeExam    int32  `json:"time_exam" binding:"required"`
}

type StartContestRequest struct {
	RId int64 `json:"r_id" binding:"required"`
}

func NewContestController(sv services.IContestService) *ContestController {
	return &ContestController{
		contestService: sv,
	}
}

func (contest *ContestController) CreateContest(ctx *gin.Context) {
	var req CreateContestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//TODO: Tùy chỉnh lỗi
		res := response.ErrorResponse(response.ErrInvalidContestSubjectName)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	contest.contestService.CreateContest(ctx, req.SubjectID, req.NumQuestion, req.TimeExam)

}
func (contest *ContestController) LiveContest(ctx *gin.Context) {
	rid := ctx.Param("id")
	parse, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrSystem))
		return
	}
	contest.contestService.LiveContest(ctx, parse)
}
func (contest *ContestController) PlayContest(ctx *gin.Context) {
	rid := ctx.Param("id")
	parse, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrSystem))
		return
	}
	contest.contestService.PlayContest(ctx, parse)
}

func (contest *ContestController) EndContest(ctx *gin.Context) {
	rid := ctx.Param("id")
	parse, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrSystem))
		return
	}
	contest.contestService.EndContest(ctx, parse)
}

func (contest *ContestController) SubmitAnswer(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Submit Contest Success",
	})
}

func (contest *ContestController) GetContestLiveById(ctx *gin.Context) {

}

func (contest *ContestController) GetListLiveContest(ctx *gin.Context) {

}

func (contest *ContestController) GetMyContestLive(ctx *gin.Context) {

}
