package entity

import "time"

type User struct {
	SessionID  string
	DequeuedAt time.Time
	ExpiredAt  time.Time
}
