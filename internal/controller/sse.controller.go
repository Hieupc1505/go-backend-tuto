package controllers

import (
	"context"
	"errors"
	"fmt"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
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

	rid := c.DefaultQuery("rid", "default")
	parseRid, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		rsp := response.ErrorResponse(response.ErrInvalidData)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	handleRoomId(c)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a new channel for this connection
	clientChan := make(chan string)
	contestID := parseRid
	global.RoomManage.AddRoom(contestID, clientChan)

	defer func() {
		global.RoomManage.RemoveRoom(contestID, clientChan)
		close(clientChan)
	}()

	// Gửi dữ liệu cho client
	for {
		select {
		case message := <-clientChan:
			if err := s.SendMessage(c.Writer, message); err != nil {
				fmt.Println("Failed to send message:", err)
				return
			}
		case <-c.Request.Context().Done():
			fmt.Println("Client disconnected")
			return
		}
	}
}

func (s *SseController) SseStartContest(c *gin.Context) {

	rid := c.DefaultQuery("rid", "default")
	parseRid, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		rsp := response.ErrorResponse(response.ErrInvalidData)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	contestID := parseRid
	message := "Start Contest"
	global.RoomManage.BroadcastToRoom(contestID, message)

	c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
}

func (s *SseController) SseEndContest(c *gin.Context) {
	message := "End Contest"
	fmt.Println(message, len(global.RoomManage.Rooms))
}

func (s *SseController) CreateContest(c *gin.Context) {

}

func (s *SseController) LiveContest(c *gin.Context) {

}

func (s *SseController) SendMessage(w http.ResponseWriter, message string) error {
	_, err := fmt.Fprintf(w, "data: %s\n\n", message)
	if err == nil {
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
	return err
}
