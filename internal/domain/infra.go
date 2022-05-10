package domain

import (
	"context"
	"time"
)

type Delivery interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type Database interface {
	LogRepository
	TelegramRepository
}

type LogRepository interface {
	SaveAppLogs(userID int, header string, body string, status int) error
}

type TelegramRepository interface {
	GetTime() (time.Time, error)
}
