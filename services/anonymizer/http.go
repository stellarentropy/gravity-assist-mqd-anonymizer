package anonymizer

import (
	"context"
	_ "expvar"
	"fmt"
	"net"
	"sync"

	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/health"

	config "github.com/stellarentropy/gravity-assist-mqd-anonymizer/config/anonymizer"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stellarentropy/gravity-assist-common/errors"
	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/config/common"
)

// getRouter returns a pointer to an instance of [*chi.Mux] that is
// pre-configured with routes and associated handler functions for the
// application. It includes routes for basic server functions, debugging, and
// application-specific endpoints with appropriate middleware applied.
func setupRouter(app *fiber.App) {
	app.All("/*", ForwardRequest)

	for _, v := range app.GetRoutes() {
		logger.Info().Str("method", v.Method).Str("path", v.Path).Msg("registered route")
	}
}

// Listen initializes an HTTP server with a predefined configuration and begins
// listening for incoming requests. On receiving a termination signal in the
// provided context, it shuts down the server gracefully. It returns an error if
// the server fails to start or encounters an issue while running.
func Listen(ctx context.Context) error {
	listener, err := net.Listen("tcp",
		net.JoinHostPort(
			config.Anonymizer.ListenAddress,
			fmt.Sprintf("%d", config.Anonymizer.ListenPort)))
	if err != nil {
		return err
	}

	// Log the starting of the anonymizer server
	logger.Info().
		Str("address", listener.Addr().(*net.TCPAddr).IP.String()).
		Int("port", listener.Addr().(*net.TCPAddr).Port).
		Msg("starting anonymizer server")

	app := fiber.New(fiber.Config{
		ReadTimeout:           common_config.Common.HealthReadTimeout,
		WriteTimeout:          common_config.Common.HealthWriteTimeout,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		StreamRequestBody:     true,
		DisableStartupMessage: true,
	})
	setupRouter(app)

	// Start a goroutine that waits for the context to be done and then shuts down the server
	go func() {
		<-ctx.Done()
		logger.Info().Msg("shutting down anonymizer server")

		if err := app.ShutdownWithTimeout(common_config.Common.GracefulShutdownTimeout); err != nil {
			// Log any errors that occur during server shutdown
			logger.Error().Err(err).Msg("error shutting down anonymizer server")
		}
	}()

	// Mark the server as ready for incoming requests
	health.Ready()

	// Start the server and return any errors that occur
	return app.Listener(listener)
}

func Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := Listen(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error().Err(err).Msg("error in anonymizer server")
			panic(err)
		}
	}()
}
