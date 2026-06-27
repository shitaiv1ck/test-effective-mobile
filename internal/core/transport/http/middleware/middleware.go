package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/logger"
	httpresponse "github.com/shitaiv1ck/test-effective-mobile/internal/core/transport/http/response"
	"go.uber.org/zap"
)

type Middleware func(next http.Handler) http.Handler

const (
	requestIDHeader = "X-Requset-ID"
	originHeader    = "Origin"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := map[string]struct{}{
				"http://localhost:8080": {},
				"null":                  {},
			}

			origin := r.Header.Get(originHeader)
			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()

				r.Header.Set(requestIDHeader, requestID)
			}

			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(logger *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.With(
				zap.String("request_id", r.Header.Get(requestIDHeader)),
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "logger", logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.FromContext(r.Context())
			rw := httpresponse.NewRW(w)

			logger.Debug("incoming http request")
			next.ServeHTTP(rw, r)
			logger.Debug(
				"done http request",
				zap.Int("status_code", rw.GetStatusCode()),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.FromContext(r.Context())

			defer func() {
				if p := recover(); p != nil {
					logger.Panic("unexpected panic", zap.Any("panic", p))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func ChainMiddleware(h http.Handler, m ...Middleware) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}
