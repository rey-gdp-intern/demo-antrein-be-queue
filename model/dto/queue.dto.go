package dto

type QueueEvent struct {
	QueueNumber   int     `json:"queue_number"`
	TimeRemaining float64 `json:"time_remaining"`
	IsFinished    bool    `json:"is_finished"`
	MainRoomToken string  `json:"main_room_token,omitempty"`
}

type RegisterQueueResponse struct {
	WaitingRoomToken string `json:"waiting_room_token"`
	MainRoomToken    string `json:"main_room_token"`
}
