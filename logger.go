package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

// newLogger returns logger with level taken from LOG_LEVEL and service from
// SERVICE_NAME environment variables.
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

	return l
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// loggerWithContext returns logger enriched with information from lambda
// context.
func loggerWithContext(ctx context.Context, logger *slog.Logger) *slog.Logger {
	lc, ok := lambdacontext.FromContext(ctx)
	if ok {
		return logger.With(
			slog.String("aws_request_id", lc.AwsRequestID),
			slog.String("invoked_function_arn", lc.InvokedFunctionArn),
		)
	}
	return logger
}
