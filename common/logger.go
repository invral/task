package common

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"time"
)

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

type loggerCtxKey struct{}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerCtxKey{}).(*slog.Logger)
	if !ok {
		return slog.New(NewDiscardHandler())
	}

	return logger
}

func FromRequest(r *http.Request) *slog.Logger {
	return FromContext(r.Context())
}

// NewHandler injects log into requests context.
func NewHandler(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(WithContext(r.Context(), log))
			next.ServeHTTP(w, r)
		})
	}
}

// AccessHandler returns a handler that call f after each request.
func AccessHandler(f func(r *http.Request, status, size int, duration time.Duration)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			f(r, ww.Status(), ww.BytesWritten(), time.Since(start))
		})
	}
}

var _ slog.Handler = &DiscardHandler{}

type DiscardHandler struct{}

func NewDiscardHandler() slog.Handler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}
