package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/onyxia-datalab/onyxia-onboarding/domain/usercontext"
	"github.com/onyxia-datalab/onyxia-onboarding/infrastructure/logging"
)

// InitLogger initializes the global logger and handles log flushing on exit.
func InitLogger(userReaderCtx usercontext.UserContextReader) {
	logger := logging.NewLogger(userReaderCtx)
	slog.SetDefault(logger)

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		slog.Info("Flushing logs before exit...")

		logging.FlushLogger()
		slog.Info("Logs successfully flushed")

		os.Exit(0)
	}()
}
