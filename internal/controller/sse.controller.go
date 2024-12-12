package controllers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type SseController struct {
	connections map[chan string]struct{}
	lock        sync.Mutex
}

func NewSseController() *SseController {
	return &SseController{
		connections: make(map[chan string]struct{}),
	}
}

func (s *SseController) SseConnection(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a new channel for this connection
	clientChan := make(chan string)

	s.lock.Lock()
	s.connections[clientChan] = struct{}{}
	fmt.Println("New connection added. Total connections:", len(s.connections))
	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		delete(s.connections, clientChan)
		s.lock.Unlock()
		close(clientChan)
	}()

	for {
		select {
		case message := <-clientChan:
			if err := s.SendMessage(c.Writer, message); err != nil {
				fmt.Println("Failed to send message:", err)
				return
			}
		case <-c.Request.Context().Done():
			fmt.Println("Client disconnected", len(s.connections))
			return
		}
	}
}

func (s *SseController) SseStartContest(c *gin.Context) {
	message := "Start Contest"
	fmt.Println(message, len(s.connections))

	// s.lock.Lock()
	// for clientChan := range s.connections {
	// 	clientChan <- message
	// }
	// s.lock.Unlock()
	s.lock.Lock()
	for clientChan := range s.connections {
		go func(ch chan string) {
			select {
			case ch <- message:
			default:
				// Skip blocked clients
			}
		}(clientChan)
	}
	s.lock.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
}

func (s *SseController) SseEndContest(c *gin.Context) {
	message := "End Contest"
	fmt.Println(message, len(s.connections))
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