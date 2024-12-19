package room

import (
	"fmt"
	"hieupc05.github/backend-server/internal/utils/sse"
	"hieupc05.github/backend-server/response"
	"sync"
)

type State struct {
	Members map[int64]chan sse.SseStatus // Map of user ID to their message channel
	State   chan bool                    // A channel to hold any state information
	Close   chan bool                    // A channel to handle remove room
}

type Manager struct {
	Rooms map[int64]*State // Map of room ID to its State
	once  sync.Once        // Ensures Rooms map initialization occurs only once
	mu    sync.RWMutex     // Mutex for synchronizing access to Rooms
}

func (m *Manager) initRooms() {
	m.once.Do(func() {
		m.Rooms = make(map[int64]*State)
	})
}

// MakeRoom creates a new room with the given ID
func (m *Manager) MakeRoom(roomID int64) {
	if m.Rooms == nil {
		m.initRooms()
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the room already exists
	if _, exists := m.Rooms[roomID]; !exists {
		m.Rooms[roomID] = &State{
			Members: make(map[int64]chan sse.SseStatus), // Initialize members
			State:   make(chan bool),                    // Initialize state channel
			Close:   make(chan bool),                    // Initialize close channel
		}
	}

}

// AddMember adds a user to the given room
func (m *Manager) AddMember(user int64, roomID int64, ch chan sse.SseStatus) int {
	if m.Rooms == nil {
		m.initRooms()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the room exists
	state, exists := m.Rooms[roomID]
	if !exists {
		return response.ErrInvalidContestGameID
	}

	// Check if the user already exists in the room
	if _, exists := state.Members[user]; exists {
		return response.ErrInvalidContestGameID
	}

	// Add the user to the room
	state.Members[user] = ch
	return 0 // Success
}

// RemoveMember removes a user from the given room
func (m *Manager) RemoveMember(roomID int64, userID int64) {
	if m.Rooms == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the room exists
	if state, exists := m.Rooms[roomID]; exists {
		// Remove the user from the room
		delete(state.Members, userID)

		// If the room has no members left, delete the room
		if len(state.Members) == 0 {
			close(state.State) // Close the state channel
			delete(m.Rooms, roomID)
		}
	}
}

// RemoveRoom removes an entire room
func (m *Manager) RemoveRoom(roomID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the room exists
	if state, exists := m.Rooms[roomID]; exists {
		// Broadcast a message to all users before removing the room
		m.BroadcastToRoom(roomID, "End contest")

		// Close the state channel and delete the room
		close(state.State)
		delete(m.Rooms, roomID)
	}
}

// BroadcastToRoom sends a message to all members of a specific room
func (m *Manager) BroadcastToRoom(roomID int64, message sse.SseStatus) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if the room exists
	if state, exists := m.Rooms[roomID]; exists {
		fmt.Printf("Broadcasting to room %d with %d members\n", roomID, len(state.Members))
		// Send the message to all members
		for userID, ch := range state.Members {
			select {
			case ch <- message:
				// Message sent successfully
				fmt.Printf("Sent message to user %d in room %d\n", userID, roomID)
			default:
				// If a user's channel is full, skip them to avoid deadlock
				fmt.Printf("Channel for user %d in room %d is full, skipping\n", userID, roomID)
			}
		}
	}
}

// IsRoomNotExist checks if a room does not exist
func (m *Manager) IsRoomNotExist(roomID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if the room exists
	_, exists := m.Rooms[roomID]
	return !exists
}
