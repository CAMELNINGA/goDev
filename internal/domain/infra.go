package domain

import (
	"context"
	"io"
	"time"
)

type Delivery interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	UploadMultipartFile(file io.ReadCloser, username string, unit string, fileName string) (string, error)
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
	GetUser(chatID int) (*User, error)
}
