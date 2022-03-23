package http

import (
	"Yaratam/internal/domain"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func generateFields(r *http.Request) logrus.Fields {
	fields := logrus.Fields{
		"ts":          time.Now().UTC().Format(time.RFC3339),
		"http_proto":  r.Proto,
		"http_method": r.Method,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"uri":         r.RequestURI,
	}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		fields["req_id"] = reqID
	}

	return fields
}

func j(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", strconv.Itoa(code))
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("cannot write response: %w", err)
	}

	return nil
}

func jError(w http.ResponseWriter, err error) error {
	code := http.StatusInternalServerError
	localizedError := "Внутренняя ошибка!"

	switch err {
	case domain.ErrInternalDatabase:
		code = http.StatusInternalServerError
		localizedError = "Внутренняя ошибка базы данных!"
	case domain.ErrInternalOpenAPI:
		code = http.StatusInternalServerError
		localizedError = "Возникла ошибка при запросе к банковским системам!"
	case domain.ErrSQLNoRows:
		code = http.StatusInternalServerError
		localizedError = "Нет данных!"
	case domain.ErrUnauthorized:
		code = http.StatusUnauthorized
		localizedError = "Вы не авторизованы!"
	case domain.ErrInvalidInputData:
		code = http.StatusBadRequest
		localizedError = "Неверный запрос!"
	case domain.ErrPaymentsLessThanAmount:
		code = http.StatusBadRequest
		localizedError = "Сумма платежей меньше суммы кредита!"
	case domain.ErrTooHighPayment:
		code = http.StatusBadRequest
		localizedError = "Слишком большой платеж!"
	case domain.ErrAmountLessThanPayment:
		code = http.StatusBadRequest
		localizedError = "Ежемесячный платеж превышает сумму кредита!"
	case domain.ErrValidationFailed:
		code = http.StatusBadRequest
		localizedError = "Запрос не прошёл валидацию!"
	case domain.ErrNoUser:
		code = http.StatusUnauthorized
		localizedError = "Вы не авторизованы!"
	case domain.ErrUnconfirmedEmail:
		code = http.StatusBadRequest
		localizedError = "Необходимо подтвердить почту!"
	case domain.ErrInvalidSMSCode:
		code = http.StatusBadRequest
		localizedError = "Неверный код!"
	case domain.ErrDuplicateRequest:
		code = http.StatusBadRequest
		localizedError = "Дублирование запроса!"
	case domain.ErrNoFCMToken:
		code = http.StatusInternalServerError
		localizedError = "Нет токена!"
	case domain.ErrIncorrectOs:
		code = http.StatusInternalServerError
		localizedError = "Нет информации об ос устройства!"
	case domain.ErrPushFailed:
		code = http.StatusInternalServerError
		localizedError = "Не получилось отправить уведомление!"
	case domain.ErrInternalPushService:
		code = http.StatusInternalServerError
		localizedError = "Ошибка на стороне push-сервиса!"
	case domain.ErrInternalIntegration:
		code = http.StatusInternalServerError
		localizedError = "Внутренняя ошибка интеграции!"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status", strconv.Itoa(code))
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error":           err.Error(),
		"localized_error": localizedError,
	}); err != nil {
		return fmt.Errorf("cannot write response: %w", err)
	}

	return nil
}
