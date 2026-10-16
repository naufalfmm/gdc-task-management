package logger

import (
	"context"
)

type LoggerData interface {
	Args() []any
}

type Logger interface {
	Info(ctx context.Context, msg string, data LoggerData)
	Error(ctx context.Context, msg string, data LoggerData)
	Debug(ctx context.Context, msg string, data LoggerData)
	Warn(ctx context.Context, msg string, data LoggerData)
}
