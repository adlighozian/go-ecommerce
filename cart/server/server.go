package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"cart-go/helper/timeout"

	"github.com/rs/zerolog"
)

func Run(srv *http.Server, logger *zerolog.Logger) error {
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
