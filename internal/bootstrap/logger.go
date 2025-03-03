package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/infrastructure/logging"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
)

// InitLogger initializes the global logger and handles log flushing on exit.
func InitLogger(userReaderCtx interfaces.UserContextReader) {
	logger, err := logging.NewLogger(userReaderCtx)

	if err != nil {
		slog.Default().
			Error("Failed to initialize logger", slog.Any("error", err))
		os.Exit(1)
	}

	slog.SetDefault(logger)

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		slog.Info("Flushing logs before exit...")

		if err := logging.FlushLogger(); err != nil {
			slog.Error("Failed to flush logs", slog.Any("error", err))
		} else {
			slog.Info("Logs successfully flushed")
		}

		os.Exit(0)
	}()
}
