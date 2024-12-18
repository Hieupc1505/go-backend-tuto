package room

import (
	"fmt"
	"sync"
)

type Manager struct {
	Rooms map[int64]map[chan string]struct{}
	once  sync.Once
	mu    sync.RWMutex
}

func (m *Manager) initRooms() {
	m.once.Do(func() {
		m.Rooms = make(map[int64]map[chan string]struct{})
	})
}

func (m *Manager) AddRoom(roomID int64, ch chan string) {
	if m.Rooms == nil {
		m.initRooms()
	}
	m.mu.Lock()         // Khóa trước khi truy cập vào Rooms
	defer m.mu.Unlock() // Mở khóa sau khi hoàn thành

	if _, exists := m.Rooms[roomID]; !exists {
		m.Rooms[roomID] = make(map[chan string]struct{})
	}
	m.Rooms[roomID][ch] = struct{}{}
}

func (m *Manager) RemoveRoom(roomID int64, ch chan string) {
	m.mu.Lock()         // Khóa trước khi truy cập vào Rooms
	defer m.mu.Unlock() // Mở khóa sau khi hoàn thành

	// Kiểm tra xem roomID có tồn tại không
	if channels, exists := m.Rooms[roomID]; exists {
		// Nếu ch == nil, xóa toàn bộ phòng
		if ch == nil {
			delete(m.Rooms, roomID)
		} else {
			// Nếu không, chỉ xóa channel cụ thể
			delete(channels, ch)
			// Nếu phòng không còn channel nào, xóa luôn phòng
			if len(channels) == 0 {
				delete(m.Rooms, roomID)
			}
		}
	}
}

func (m *Manager) BroadcastToRoom(roomID int64, message string) {
	m.mu.RLock()         // Khóa đọc để đảm bảo đồng bộ
	defer m.mu.RUnlock() // Mở khóa sau khi hoàn thành

	// Lấy danh sách các channel trong roomID
	if channels, exists := m.Rooms[roomID]; exists {
		fmt.Println("Sending message to room", roomID)
		for ch := range channels {
			// Gửi message cho từng channel
			select {
			case ch <- message:
				// Message sent successfully
			default:
				// Nếu channel bị đầy, bỏ qua để tránh deadlock
				fmt.Println("Channel full, skipping message")
			}
		}
	}
}
