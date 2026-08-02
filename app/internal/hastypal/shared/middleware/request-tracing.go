package middleware

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

func NewRequestTracingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		id := uuid.New()

		r.Header.Add("traceId", id.String())

		initialTrace := fmt.Sprintf(
			"Request: Method=%s URI=%s UserAgent=%s TraceId=%s",
			r.Method,
			r.RequestURI,
			r.UserAgent(),
			id.String(),
		)

		slog.Info(initialTrace)

		rww := helper.NewResponseWriterWrapper(w)

		handler.ServeHTTP(rww, r)

		elapsed := time.Since(start)

		trace := fmt.Sprintf(
			"Response: Method=%s URI=%s UserAgent=%s Time=%dms Status=%d TraceId=%s",
			r.Method,
			r.RequestURI,
			r.UserAgent(),
			elapsed.Milliseconds(),
			rww.StatusCode(),
			id.String(),
		)

		slog.Info(trace)
	})
}
