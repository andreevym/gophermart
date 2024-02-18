package middleware

import (
	"io"
	"net/http"
	"time"

	"github.com/andreevym/gophermart/pkg/logger"
	"go.uber.org/zap"
)

func WithRequestLoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		h.ServeHTTP(w, r)
		end := time.Now()

		if r == nil {
			return
		}

		logger.Logger().Debug(
			"request",
			zap.String("method", r.Method),
			zap.String("URI", r.RequestURI),
			zap.Duration("duration", end.Sub(start)),
		)

		if r.Response == nil {
			return
		}
		defer func() {
			err := r.Body.Close()
			if err != nil {
				logger.Logger().Error(err.Error())
			}
		}()
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger().Error("read body error", zap.Error(err))
			return
		}

		logger.Logger().Debug(
			"response",
			zap.Int("status", r.Response.StatusCode),
			zap.Int("status", len(bytes)),
		)
	})
}
