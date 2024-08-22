package entity

import "time"

type Session struct {
	SessionID  string
	EnqueuedAt time.Time
}
