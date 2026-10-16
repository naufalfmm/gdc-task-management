package slog

import (
	"context"
	"log/slog"
	"os"

	"github.com/naufalfmm/gdc-task-management/internal/pkg/logger"
)

type slogLogger struct {
	log *slog.Logger
}

func New() slogLogger {
	return slogLogger{
		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (s slogLogger) Info(ctx context.Context, msg string, data logger.LoggerData) {
	s.log.InfoContext(ctx, msg, data.Args()...)
}

func (s slogLogger) Debug(ctx context.Context, msg string, data logger.LoggerData) {
	s.log.DebugContext(ctx, msg, data.Args()...)
}

func (s slogLogger) Warn(ctx context.Context, msg string, data logger.LoggerData) {
	s.log.WarnContext(ctx, msg, data.Args()...)
}

func (s slogLogger) Error(ctx context.Context, msg string, data logger.LoggerData) {
	s.log.ErrorContext(ctx, msg, data.Args()...)
}
