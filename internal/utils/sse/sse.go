package sse

type SseStatus string // status

const (
	Disconnected SseStatus = "disconnected"
	Connected    SseStatus = "connected"
	UserJoin     SseStatus = "user_join"
	UserLeave    SseStatus = "user_leave"
	LiveContest  SseStatus = "live"
	StartContest SseStatus = "start_contest"
	EndContest   SseStatus = "end_contest"
	CloseContest SseStatus = "close_contest"
	ContestInfo  SseStatus = "contest_info"
)
