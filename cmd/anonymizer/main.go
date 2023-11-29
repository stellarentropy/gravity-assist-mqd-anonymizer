package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/stellarentropy/gravity-assist-common/logging"
	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/components"
	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/health"
	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/services/anonymizer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	logger := logging.GetLogger()

	logger.Info().Msg("starting anonymizer")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signalChan
		logger.Info().Msg("shutting down anonymizer")
		cancel()
	}()

	components.StartComponentsAndWait(ctx,
		health.NewComponent(),
		anonymizer.NewComponent(),
	)

	logger.Info().Msg("anonymizer shutdown complete")
}
