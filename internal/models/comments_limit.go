package models

import "time"

type CommentLimit struct {
	MaxAttempts int           `json:"max_attempts"`
	Window      time.Duration `json:"window"`
}
