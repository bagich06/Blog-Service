package models

import "time"

type LoginLimit struct {
	MaxAttempts int           `json:"max_attempts"`
	Window      time.Duration `json:"window"`
}
