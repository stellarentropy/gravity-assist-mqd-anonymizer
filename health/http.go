package health

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/stellarentropy/gravity-assist-common/errors"
	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/config/common"

	"github.com/gofiber/fiber/v2"

	"github.com/goccy/go-json"
)

var ready = false

func Ready() {
	ready = true
}

func setupRouter(app *fiber.App) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/readyz", func(c *fiber.Ctx) error {
		if ready {
			return c.SendString("OK")
		}
		return c.SendStatus(fiber.StatusServiceUnavailable)
	})

	for _, v := range app.GetRoutes() {
		logger.Info().Str("method", v.Method).Str("path", v.Path).Msg("registered route")
	}
}

func Listen(ctx context.Context) error {
	listener, err := net.Listen("tcp",
		net.JoinHostPort(
			common_config.Common.HealthListenAddress,
			fmt.Sprintf("%d", common_config.Common.HealthListenPort)))
	if err != nil {
		return err
	}

	// Log the starting of the health server
	logger.Info().
		Str("address", listener.Addr().(*net.TCPAddr).IP.String()).
		Int("port", listener.Addr().(*net.TCPAddr).Port).
		Msg("starting health server")

	app := fiber.New(fiber.Config{
		ReadTimeout:           common_config.Common.HealthReadTimeout,
		WriteTimeout:          common_config.Common.HealthWriteTimeout,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})
	setupRouter(app)

	// Start a new goroutine that listens for the context cancellation signal
	go func() {
		// Wait for the context to be done
		<-ctx.Done()
		// Log the shutting down of the health server
		logger.Info().Msg("shutting down health server")

		// Attempt to gracefully shut down the server and log any errors
		if err := app.ShutdownWithTimeout(common_config.Common.GracefulShutdownTimeout); err != nil {
			logger.Error().Err(err).Msg("error shutting down health server")
		}
	}()

	// Start the server and return any errors encountered
	return app.Listener(listener)
}

func Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := Listen(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error().Err(err).Msg("error in health server")
			panic(err)
		}
	}()
}
