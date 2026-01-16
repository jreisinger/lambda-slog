package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func newLogger() *slog.Logger {
	level := slog.LevelInfo
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN", "WARNING":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		// AddSource: true, // include file:line
	})

	l := slog.New(h).With(
		slog.String("service", getenv("SERVICE_NAME", "my-lambda")),
	)

	// Make it default logger so slog.Info and others work too.
	slog.SetDefault(l)

	return l
}

func main() {
	logger := newLogger()
	lambda.Start(newHandler(logger))
}

func newHandler(logger *slog.Logger) func(context.Context) (string, error) {
	var coldStart atomic.Bool
	coldStart.Store(true)

	return func(ctx context.Context) (string, error) {
		logger := withInvocationLogger(ctx, logger)

		// Cold start flag is useful in dashboards
		logger = logger.With(slog.Bool("cold_start", coldStart.Load()))
		coldStart.Store(false)

		start := time.Now()
		logger.InfoContext(ctx, "handling request")

		// Example: return an error sometimes
		if time.Now().Unix()%7 == 0 {
			err := errors.New("boom")
			logger.ErrorContext(ctx, "request failed", slog.Any("err", err))
			return "", err
		}

		logger.InfoContext(ctx, "request ok", slog.Int64("duration_ns", time.Since(start).Nanoseconds()))
		return "ok", nil
	}
}

func withInvocationLogger(ctx context.Context, logger *slog.Logger) *slog.Logger {
	lc, ok := lambdacontext.FromContext(ctx)
	if ok {
		return logger.With(
			slog.String("aws_request_id", lc.AwsRequestID),
			slog.String("invoked_function_arn", lc.InvokedFunctionArn),
		)
	}
	return logger
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
