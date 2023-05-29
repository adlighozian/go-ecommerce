package server

import (
	"api-gateway-go/helper/timeout"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

func Run(srv *http.Server, logger *zerolog.Logger) error {
	// note:
	// Graceful restarts or stops are important
	// allow us to safely shut down the server while it is still processing requests.
	// This ensures that no ongoing requests are abruptly terminated,
	// preventing data inconsistencies or loss,
	// and the server can complete the processing of any requests in progress before shutting down.
	// credit: https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	errChan := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("ListenAndServe(): %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-quit:
	}
	logger.Debug().Msg("shutdown server...")

	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	logger.Debug().Msg("server exiting")

	return nil
}
