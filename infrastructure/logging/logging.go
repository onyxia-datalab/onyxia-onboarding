package logging

import (
	"log"
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func NewLogger() *slog.Logger {
	zapLogger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}

	handler := zapslog.NewHandler(
		zapLogger.Core(),
		zapslog.WithCaller(true),
		//zapslog.AddStacktraceAt(slog.LevelDebug),
	)

	return slog.New(handler)
}

func FlushLogger() {
	logger := zap.L() // Get global zap logger
	if err := logger.Sync(); err != nil {
		slog.Error("Failed to flush logs", slog.Any("error", err))
	}
}
