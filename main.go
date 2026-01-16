package main

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	logger := newLogger()
	lambda.Start(newHandler(logger))
}

// newHandler creates a new Lambda handler function that logs request processing
// details. It sometimes returns an error for demonstration purposes.
func newHandler(logger *slog.Logger) func(context.Context) (string, error) {
	var coldStart atomic.Bool
	coldStart.Store(true)

	return func(ctx context.Context) (string, error) {
		logger := loggerWithContext(ctx, logger)

		// Cold start flag is useful in dashboards
		logger = logger.With(slog.Bool("cold_start", coldStart.Load()))
		coldStart.Store(false)

		start := time.Now()
		logger.InfoContext(ctx, "handling request")

		// Example: return an error sometimes
		if time.Now().Unix()%2 == 0 {
			err := errors.New("boom")
			logger.ErrorContext(ctx, "request failed", slog.Any("err", err))
			return "", err
		}

		logger.InfoContext(ctx, "request ok", slog.Int64("duration_ns", time.Since(start).Nanoseconds()))
		return "ok", nil
	}
}
