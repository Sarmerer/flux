package middleware

import (
	"context"
	"fmt"
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

			verbosity := getVerbosity(ctxLogger)

			if verbosity >= logging.VerbosityVerbose {
				ctxLogger.Info("HTTP request started", map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"remote_addr": r.RemoteAddr,
					"user_agent":  r.UserAgent(),
				})
			}

			defer func() {
				duration := time.Since(start)
				status := ww.Status()

				methodIcon := getMethodIcon(r.Method)
				statusIcon := getStatusIcon(status)

				if verbosity == logging.VerbosityMinimal {
					ctxLogger.Info(formatCompactLog(methodIcon, statusIcon, r.Method, r.URL.Path, status, duration))
				} else {
					fields := map[string]interface{}{
						"method": r.Method,
						"status": status,
						"ms":     duration.Milliseconds(),
					}

					if verbosity >= logging.VerbosityVerbose {
						fields["path"] = r.URL.Path
						fields["bytes"] = ww.BytesWritten()
						fields["ip"] = r.RemoteAddr
					}

					msg := fmt.Sprintf("%s %s %d", methodIcon, r.URL.Path, status)

					if status >= 500 {
						ctxLogger.Error(msg, nil, fields)
					} else if status >= 400 {
						ctxLogger.Warn(msg, fields)
					} else {
						ctxLogger.Info(msg, fields)
					}
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

func getVerbosity(logger *logging.ContextLogger) logging.LogVerbosity {
	return logger.GetVerbosity()
}

func getMethodIcon(method string) string {
	switch method {
	case "GET":
		return "↓"
	case "POST":
		return "↑"
	case "PUT":
		return "⟳"
	case "PATCH":
		return "✎"
	case "DELETE":
		return "✕"
	case "OPTIONS":
		return "⚙"
	default:
		return "•"
	}
}

func getStatusIcon(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "✓"
	case status >= 300 && status < 400:
		return "➜"
	case status >= 400 && status < 500:
		return "!"
	case status >= 500:
		return "✗"
	default:
		return "?"
	}
}

func formatCompactLog(methodIcon, statusIcon, method, path string, status int, duration time.Duration) string {
	return fmt.Sprintf("%s %s %s %d %dms", methodIcon, statusIcon, method, status, duration.Milliseconds())
}
