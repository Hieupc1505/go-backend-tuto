package controllers

import (
	"context"
	"errors"
	"fmt"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/middlewares"
	"hieupc05.github/backend-server/internal/utils/sse"
	"hieupc05.github/backend-server/internal/utils/token"
	"hieupc05.github/backend-server/response"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type SseController struct {
	lock sync.Mutex
}

func NewSseController() *SseController {
	return &SseController{}
}

func handleRoomId(c *gin.Context) {
	rid := c.DefaultQuery("rid", "default")
	parseRid, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		rsp := response.ErrorResponse(response.ErrInvalidData)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	_, err = global.PgDb.GetContest(context.Background(), parseRid)
	if err != nil {
		if errors.Is(db.ErrRecordNotFound, err) {
			rsp := response.ErrorResponse(response.ErrInvalidContestNotFound)
			c.JSON(http.StatusBadRequest, rsp)
			return
		}
		rsp := response.ErrorResponse(response.ErrSystem)
		c.JSON(http.StatusInternalServerError, rsp)
		return
	}
}

func (s *SseController) SseConnection(c *gin.Context) {

	rid := c.DefaultQuery("rid", "-1")
	parseRid, err := strconv.ParseInt(rid, 10, 64)

	if err != nil || parseRid < 0 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrInvalidContestGameID))
		return
	}

	if roomNotExists := global.RoomManage.IsRoomNotExist(parseRid); roomNotExists {
		fmt.Println("Room does not exist", roomNotExists)
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrInvalidContestGameID))
		return
	}

	// Validate room by contest ID
	handleRoomId(c)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a new channel for this connection
	authPayload, exists := c.Get(middlewares.AuthorizationPayloadKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(response.ErrInvalidContestGameID))
		return
	}

	tokenPayload, ok := authPayload.(*token.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrSystem))
		return
	}

	clientChan := make(chan sse.SseStatus)
	contestID := parseRid
	errCode := global.RoomManage.AddMember(tokenPayload.UserID, contestID, clientChan)
	if errCode != 0 {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(errCode))
		return
	}

	defer func() {
		global.RoomManage.RemoveMember(contestID, tokenPayload.UserID)
		close(clientChan)
	}()

	handleMessage(sse.UserJoin)

	// Send data to client
	for {
		select {
		case message := <-clientChan:
			rsp := handleMessage(message)
			if err := s.SendMessage(c.Writer, rsp); err != nil {
				fmt.Println("Failed to send message:", err)
				return
			}
		case <-c.Request.Context().Done():
			fmt.Println("Client disconnected")
			return
		case <-global.RoomManage.Rooms[contestID].State:
			fmt.Println("Contest ended, closing connection")
			return
		}
	}
}

func handleMessage(msg sse.SseStatus) gin.H {
	switch msg {
	case sse.UserJoin:
		return gin.H{"status": "User joined"}
	case sse.LiveContest:
		return gin.H{"status": "Live Contest"}
	case sse.EndContest:
		return gin.H{"status": "End Contest"}
	case sse.StartContest:
		return gin.H{"status": "Start Contest"}
	case sse.UserLeave:
		return gin.H{"status": "User left"}
	case sse.ContestInfo:
		return gin.H{"status": "Contest information"}
	default:
		return gin.H{"status": "Message sent"}
	}
}

func (s *SseController) SseStartContest(c *gin.Context) {

	rid := c.DefaultQuery("rid", "-1")
	parseRid, err := strconv.ParseInt(rid, 10, 64)
	if err != nil || parseRid == -1 {
		rsp := response.ErrorResponse(response.ErrInvalidData)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	contestID := parseRid
	message := sse.StartContest
	global.RoomManage.BroadcastToRoom(contestID, message)

	c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
}

func (s *SseController) SseEndContest(c *gin.Context) {
	message := "End Contest"
	fmt.Println(message, len(global.RoomManage.Rooms))
}

func (s *SseController) SendMessage(w http.ResponseWriter, message gin.H) error {
	_, err := fmt.Fprintf(w, "data: %s\n\n", message)
	if err == nil {
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
	return err
}
