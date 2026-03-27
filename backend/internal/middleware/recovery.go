package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/rastignacc/recircuit-demo-app/backend/internal/model"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						slog.Any("error", err),
						slog.String("stack", string(debug.Stack())),
						slog.String("request_id", GetRequestID(r.Context())),
					)
					model.WriteError(w, model.ErrInternal("internal server error"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
