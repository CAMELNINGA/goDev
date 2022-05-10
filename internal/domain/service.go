package domain

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	LogService
	TelegramService
}

type LogService interface {
	SaveAppLogs(w http.ResponseWriter, r *http.Request, next http.Handler) error
}

type TelegramService interface {
	GetTimes() (time.Time, error)
}

type service struct {
	logger logrus.FieldLogger
	db     Database
}

func NewService(logger logrus.FieldLogger, db Database) Service {
	s := &service{
		logger: logger,
		db:     db,
	}
	return s
}

func (s *service) SaveAppLogs(w http.ResponseWriter, r *http.Request, next http.Handler) error {
	userID := r.Context().Value(ContextUserID).(int)

	header := r.URL.Path
	method := r.Method

	next.ServeHTTP(w, r)

	statusSTR, err := strconv.Atoi(w.Header().Get("Status"))
	if err != nil {
		s.logger.WithField("user_id", userID).WithError(err).Error("status dosent converting to int")
	}
	return s.db.SaveAppLogs(userID, header, method, statusSTR)
}

func (s *service) GetTimes() (time.Time, error) {
	return s.db.GetTime()
}
