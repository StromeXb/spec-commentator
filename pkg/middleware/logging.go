package middleware

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog"
)

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger *zerolog.Logger, endpointName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logger.Info().Msg(fmt.Sprintf("Call %s endpoint", endpointName))
			response, err = next(ctx, request)
			if err != nil {
				logger.Error().Err(err).Msg(fmt.Sprintf("Error call %s endpoint", endpointName))
			}
			return response, err
		}
	}
}

func GetLoggingMiddlewareFunc(logger *zerolog.Logger) func(field reflect.StructField) endpoint.Middleware {
	return func(field reflect.StructField) endpoint.Middleware {
		return LoggingMiddleware(logger, field.Name)
	}
}

func SetLoggingMiddleware(loggingProvider LoggingProvider, endpoints interface{}) {
	mwFunc := GetLoggingMiddlewareFunc(loggingProvider.Logger())
	SetMiddleware(mwFunc, reflect.ValueOf(endpoints).Elem())
}

type LoggingProvider interface {
	Logger() *zerolog.Logger
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// LoggingTransportMiddleware logs the incoming HTTP request & its duration.
func LoggingTransportMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error().Str("method", r.Method).
						Str("path", r.URL.EscapedPath()).
						Str("query", r.URL.RawQuery).
						Str("trace", string(debug.Stack())).
						Interface("err", err).
						Msg("Error when calling HTTP method")
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Debug().
				Int("status", wrapped.status).
				Str("method", r.Method).
				Str("path", r.URL.EscapedPath()).
				Str("query", r.URL.RawQuery).
				Dur("duration", time.Since(start)).
				Msg("HTTP method called")
		}

		return http.HandlerFunc(fn)
	}
}
