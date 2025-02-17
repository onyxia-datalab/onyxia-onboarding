package logging

import (
	"log"
	"log/slog"

	"github.com/onyxia-datalab/onyxia-onboarding/domain/usercontext"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func NewLogger(userReaderCtx usercontext.UserContextReader) *slog.Logger {
	zapLogger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}

	baseHandler := zapslog.NewHandler(
		zapLogger.Core(),
		zapslog.WithCaller(true),
	)

	contextHandler := NewUserContextLogger(baseHandler, userReaderCtx)

	return slog.New(contextHandler)
}

func FlushLogger() {
	logger := zap.L() // Get global zap logger
	if err := logger.Sync(); err != nil {
		slog.Error("Failed to flush logs", slog.Any("error", err))
	}
}
