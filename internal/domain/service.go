package domain

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type Service interface {
	LogService
	TelegramService
}

type LogService interface {
	SaveAppLogs(w http.ResponseWriter, r *http.Request, next http.Handler) error
}

type TelegramService interface {
	GetUserData(ID int) (*User, error)
	AddUser(user *User) error
	UploadMultipartFile(file io.ReadCloser, username string, unit string, fileName string) (string, error)
}

type service struct {
	logger logrus.FieldLogger
	db     Database
	http   Httperf
}

func NewService(logger logrus.FieldLogger, db Database, httpreq Httperf) Service {
	s := &service{
		logger: logger,
		db:     db,
		http:   httpreq,
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

func (s *service) GetUserData(ID int) (*User, error) {
	return s.db.GetUser(ID)
}

func (s *service) AddUser(user *User) error {
	return s.db.AddUser(user)
}

func (s *service) UploadMultipartFile(file io.ReadCloser, username string, unit string, fileName string) (string, error) {
	return s.http.UploadMultipartFile(file, username, unit, fileName)
}
