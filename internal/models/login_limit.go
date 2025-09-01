package models

import "time"

type LoginLimit struct {
	MaxAttempts int           `json:"max_attempts"`
	Window      time.Duration `json:"window"` // в течении какого времени были выполнены попытки
}
