package domain

import "context"

type Delivery interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type Database interface {
	LogRepository
}

type LogRepository interface {
	SaveAppLogs(userID int, header string, body string, status int) error
}
