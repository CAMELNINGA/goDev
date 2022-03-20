package http

import (
	"Yaratam/internal/domain"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
	"time"
)

type appLogger struct {
	logger logrus.FieldLogger
}

func (al *appLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
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

	logger := al.logger.WithFields(fields)
	logger.Infoln("request started")

	return &appLoggerEntry{
		Logger: logger,
	}
}

type appLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (a appLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, _ interface{}) {
	a.Logger = a.Logger.WithFields(logrus.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed_ms":   float64(elapsed.Nanoseconds()) / 1_000_000.0,
	})

	if userID := header.Get("User-ID"); userID != "" {
		a.Logger = a.Logger.WithField("user_id", userID)
	}

	a.Logger.Infoln("request complete")
}

func (a appLoggerEntry) Panic(v interface{}, stack []byte) {
	a.Logger = a.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func (a *adapter) loggerMiddleware(next http.Handler) http.Handler {
	logger := &appLogger{logger: a.logger}

	fn := func(w http.ResponseWriter, r *http.Request) {
		entry := logger.NewLogEntry(r)
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t := time.Now()
		defer func() {
			entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t), nil)
		}()

		next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
	}

	return http.HandlerFunc(fn)
}

func (a *adapter) logSavingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_ = a.service.SaveAppLogs(w, r, next)
	}
	return http.HandlerFunc(fn)
}

func (a *adapter) JWTAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authMiddleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					log.Println("middleware")
					// authorizationHeader := r.Header.Get("Authorization")
					// if strings.HasPrefix(authorizationHeader, "Bearer") {
					//   r.Header.Set("Authorization", jwtauth.TokenFromHeader(r))
					// }

					// token, err := jwtauth.VerifyRequest(a.jwtAuth, r, func(r *http.Request) string {
					//   return r.Header.Get("Authorization")
					// })

					// if err != nil {
					//   a.logger.WithError(err).Error(domain.ErrUnauthorized)
					//   _ = jError(w, domain.ErrUnauthorized)
					//   return
					// }

					ctx := r.Context()
					// ctx = jwtauth.NewContext(ctx, token, err)
					next.ServeHTTP(w, r.WithContext(ctx))
				})

			}

			userIDMiddleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					//     _, claims, err := jwtauth.FromContext(r.Context())
					//     if err != nil {
					//       a.logger.WithError(err).Error("Error while extracting JWT token from context!")
					//       _ = jError(w, domain.ErrUnauthorized)
					//       return
					//     }

					//     sub, ok := claims["sub"]
					//     if !ok {
					//       a.logger.Error("Token is without 'sub' field!")
					//       _ = jError(w, domain.ErrUnauthorized)
					//       return
					//     }

					//     subStr, ok := sub.(string)
					//     if !ok {
					//       a.logger.Error("'sub' field is not a string!")
					//       _ = jError(w, domain.ErrUnauthorized)
					//       return
					//     }

					// id, err := strconv.Atoi(subStr)
					// if err != nil {
					//   a.logger.WithError(err).Error("Error while converting string into int!")
					//   _ = jError(w, domain.ErrUnauthorized)
					//   return
					// }

					r = r.WithContext(context.WithValue(r.Context(), domain.ContextUserID, 1))
					w.Header().Set("User-ID", strconv.Itoa(1))
					next.ServeHTTP(w, r)

				})
			}
			chi.Chain(authMiddleware, userIDMiddleware).Handler(next).ServeHTTP(w, r)
		})
	}
}
