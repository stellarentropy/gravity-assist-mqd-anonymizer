package anonymizer

import (
	"context"
	_ "expvar"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/health"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	config "github.com/stellarentropy/gravity-assist-mqd-anonymizer/config/anonymizer"

	"github.com/go-chi/chi/v5"
	"github.com/stellarentropy/gravity-assist-common/errors"
	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/config/common"
)

// getRouter returns a pointer to an instance of [*chi.Mux] that is
// pre-configured with routes and associated handler functions for the
// application. It includes routes for basic server functions, debugging, and
// application-specific endpoints with appropriate middleware applied.
func getRouter() *chi.Mux {
	router := chi.NewRouter()

	router.HandleFunc("/", ForwardRequest)

	for _, v := range router.Routes() {
		logger.Info().Str("pattern", v.Pattern).Msg("registered route")
	}

	return router
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

	// Get the address of the server
	addr := listener.Addr().(*net.TCPAddr).String()

	router := getRouter()

	server := &http.Server{
		Addr:         addr,
		Handler:      h2c.NewHandler(router, &http2.Server{}),
		ReadTimeout:  common_config.Common.HealthReadTimeout,
		WriteTimeout: common_config.Common.HealthWriteTimeout,
	}

	// Start a goroutine that waits for the context to be done and then shuts down the server
	go func() {
		<-ctx.Done()
		logger.Info().Msg("shutting down anonymizer server")

		tctx, cancel := context.WithTimeout(context.Background(), common_config.Common.GracefulShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(tctx); err != nil {
			// Log any errors that occur during server shutdown
			logger.Error().Err(err).Msg("error shutting down anonymizer server")
		}
	}()

	// Mark the server as ready for incoming requests
	health.Ready()

	// Start the server and return any errors that occur
	return server.Serve(listener)
}

func Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := Listen(ctx); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("error in anonymizer server")
			panic(err)
		}
	}()
}
