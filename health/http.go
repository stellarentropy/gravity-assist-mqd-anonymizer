package health

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/stellarentropy/gravity-assist-common/errors"
	"github.com/stellarentropy/gravity-assist-mqd-anonymizer/config/common"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var ready = false

func Ready() {
	ready = true
}

func getRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if ready {
			_, _ = w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})

	for _, v := range router.Routes() {
		logger.Info().Str("pattern", v.Pattern).Msg("registered route")
	}

	return router
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

	router := getRouter()

	// Get the address of the server
	addr := listener.Addr().(*net.TCPAddr).String()

	server := &http.Server{
		Addr:         addr,
		Handler:      h2c.NewHandler(router, &http2.Server{}),
		ReadTimeout:  common_config.Common.HealthReadTimeout,
		WriteTimeout: common_config.Common.HealthWriteTimeout,
	}

	// Start a new goroutine that listens for the context cancellation signal
	go func() {
		// Wait for the context to be done
		<-ctx.Done()
		// Log the shutting down of the health server
		logger.Info().Msg("shutting down health server")

		tctx, cancel := context.WithTimeout(context.Background(), common_config.Common.GracefulShutdownTimeout)
		defer cancel()

		// Attempt to gracefully shut down the server and log any errors
		if err := server.Shutdown(tctx); err != nil {
			logger.Error().Err(err).Msg("error shutting down health server")
		}
	}()

	// Start the server and return any errors encountered
	return server.Serve(listener)
}

func Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := Listen(ctx); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("error in health server")
			panic(err)
		}
	}()
}
