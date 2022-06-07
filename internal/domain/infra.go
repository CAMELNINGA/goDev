package domain

import (
	"context"
	"io"
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
	GetUser(chatID int) (*User, error)
	AddUser(user *User) error
	AddPath(chatID int, path *Path) (int, error)
	ChangeUserPath(chatID, pathID int) error
	AddFile(chatID int, path string) error
	GetFiles(chatID int64) ([]*File, error)
	GetPaths(chatID int) ([]*Path, error)
	DeletePaths(chatID int) error
}

type Httperf interface {
	UploadMultipartFile(file io.ReadCloser, username string, unit string, fileName string) (string, error)
}
