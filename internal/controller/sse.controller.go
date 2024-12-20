package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "hieupc05.github/backend-server/db/sqlc"
	"hieupc05.github/backend-server/global"
	"hieupc05.github/backend-server/internal/middlewares"
	"hieupc05.github/backend-server/internal/utils/sse"
	"hieupc05.github/backend-server/internal/utils/token"
	"hieupc05.github/backend-server/response"
	"net/http"
	"strconv"
	"sync"
	"time"
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

	clientChan := make(chan sse.SseStatus, 3)
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

	// Sequentially broadcast UserJoin and ContestInfo
	err = global.RoomManage.BroadcastToRoom(contestID, sse.ContestInfo, sse.UserJoin)
	if err != nil {
		fmt.Println("Failed to broadcast UserJoin message:", err)
		return
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()
	c.Request = c.Request.WithContext(ctx)

	// Send data to client
	for {
		select {
		case message := <-clientChan:
			fmt.Println("Message received:", message)
			data, err := handleMessage(contestID, message)
			if err != nil {
				fmt.Println("Failed to handle message:", err)
				rsp := SseResponse{
					Error: 1,
					Data:  SseData{Code: "message", Data: "Failed to handle message"},
				}
				if err := s.SendMessage(c.Writer, rsp); err != nil {
					fmt.Println("Failed to send message:", err)
					return
				}
			}
			rsp := SseResponse{
				Error: 0,
				Data:  SseData{Code: "message", Data: data},
			}
			if err := s.SendMessage(c.Writer, rsp); err != nil {
				fmt.Println("Failed to send message:", err)
				return
			}
			if message == sse.ContestResults {
				cancel()
				fmt.Println("Remove room")
			}

		case <-c.Request.Context().Done():
			global.RoomManage.RemoveMember(contestID, tokenPayload.UserID)
			global.RoomManage.RemoveRoom(contestID)
			fmt.Println("Client disconnected")
			return
		case <-global.RoomManage.Rooms[contestID].State:
			fmt.Println("Contest ended, closing connection, handleResult")
			time.Sleep(5 * time.Second)
			err := global.RoomManage.BroadcastToRoom(contestID, sse.ContestResults)
			if err != nil {
				fmt.Println("Failed to broadcast ContestResults message:", err)
			}

		}
	}
}

type SseData struct {
	Code string      `json:"pkg_code"`
	Data interface{} `json:"pkg_data"`
}

type SseResponse struct {
	Error int     `json:"e"`
	Data  SseData `json:"d"`
}

func handleMessage(roomID int64, msg sse.SseStatus) (interface{}, error) {
	switch msg {
	case sse.UserJoin:
		fmt.Println("User joined")
		return []map[string]interface{}{
			{
				"id":       184,
				"owner":    true,
				"nickname": "Taka IT",
			},
			{
				"id":       235,
				"owner":    false,
				"nickname": "Taka IT",
				"results": map[string]interface{}{
					"num_correct":   1,
					"num_incorrect": 1,
					"time_submit":   140,
					"results": []map[string]interface{}{
						{
							"question_id": 96,
							"exam": map[string]interface{}{
								"index":      0,
								"ans":        "ans 2",
								"is_correct": true,
							},
							"correct": map[string]interface{}{
								"index":      0,
								"ans":        "ans 2",
								"is_correct": true,
							},
						},
						{
							"question_id": 97,
							"exam": map[string]interface{}{
								"index":      0,
								"ans":        "ans 1",
								"is_correct": false,
							},
							"correct": map[string]interface{}{
								"index":      0,
								"ans":        "ans 2",
								"is_correct": true,
							},
						},
					},
				},
			},
		}, nil
	case sse.LiveContest:
		return gin.H{"status": "Live Contest"}, nil
	case sse.EndContest:
		if global.RoomManage.IsRoomNotExist(roomID) {
			return gin.H{"status": "Room does not exist"}, nil
		}
		return gin.H{"status": "End Contest"}, nil
	case sse.StartContest:
		return gin.H{"status": "Start Contest"}, nil
	case sse.UserLeave:
		return gin.H{"status": "User left"}, nil
	case sse.ContestInfo:
		fmt.Println("Contest info")
		contest, err := global.PgDb.GetContest(context.Background(), roomID)
		if err != nil {
			return "", err
		}
		return contest, nil
	case sse.ContestResults:
		return gin.H{"status": "Contest results"}, nil
	default:
		return gin.H{"status": "Message sent"}, nil
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
	err = global.RoomManage.BroadcastToRoom(contestID, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrInvalidContestGameID))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
}

func (s *SseController) SseEndContest(c *gin.Context) {
	message := "End Contest"
	fmt.Println(message, len(global.RoomManage.Rooms))
}

func (s *SseController) SendMessage(w http.ResponseWriter, data SseResponse) error {
	// Encode data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write JSON to response (without extra strings)
	_, err = fmt.Fprintf(w, "%s\n\n", jsonData)
	if err == nil {
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
	return err
}
