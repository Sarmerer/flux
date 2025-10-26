package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/flow/internal/infrastructure/logging"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

const LoggerKey contextKey = "logger"

func LoggingMiddleware(appLogger *logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := middleware.GetReqID(r.Context())
			if requestID == "" {
				requestID = uuid.New().String()
			}

			logCtx := logging.LogContext{
				RequestID: requestID,
				Operation: r.Method + " " + r.URL.Path,
			}

			userID := r.Context().Value("user_id")
			if userID != nil {
				if uid, ok := userID.(uuid.UUID); ok {
					logCtx.UserID = &uid
				}
			}

			ctxLogger := appLogger.WithContext(logCtx)

			ctx := context.WithValue(r.Context(), LoggerKey, ctxLogger)
			r = r.WithContext(ctx)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			ctxLogger.Info("HTTP request started", map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
			})

			defer func() {
				duration := time.Since(start)
				status := ww.Status()

				fields := map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"status":      status,
					"bytes":       ww.BytesWritten(),
					"duration_ms": duration.Milliseconds(),
					"remote_addr": r.RemoteAddr,
				}

				if status >= 500 {
					ctxLogger.Error("HTTP request failed", nil, fields)
				} else if status >= 400 {
					ctxLogger.Warn("HTTP request completed with client error", fields)
				} else {
					ctxLogger.Info("HTTP request completed", fields)
				}
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func GetLogger(ctx context.Context) *logging.ContextLogger {
	if logger, ok := ctx.Value(LoggerKey).(*logging.ContextLogger); ok {
		return logger
	}
	return nil
}
