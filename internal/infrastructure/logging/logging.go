package logging

import (
	"fmt"
	"log/slog"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func NewLogger(userReaderCtx interfaces.UserContextReader) (*slog.Logger, error) {
	zapLogger, err := zap.NewProduction()

	if err != nil {
		return nil, fmt.Errorf("failed to initialize zap logger: %w", err)

	}

	baseHandler := zapslog.NewHandler(
		zapLogger.Core(),
		zapslog.WithCaller(true),
	)

	contextHandler := NewUserContextLogger(baseHandler, userReaderCtx)

	return slog.New(contextHandler), nil
}

func FlushLogger() error {
	logger := zap.L() // Get global zap logger
	if err := logger.Sync(); err != nil {
		slog.Error("Failed to flush logs", slog.Any("error", err))
		return err
	}
	return nil
}
