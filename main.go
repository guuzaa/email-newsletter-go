package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guuzaa/email-newsletter/cmd"
	"github.com/guuzaa/email-newsletter/internal"
)

var logger = internal.Logger()

func main() {
	config, err := internal.Configuration("configuration")
	if err != nil {
		logger.Panic().Err(err)
	}

	srv, err := cmd.Build(&config)
	if err != nil {
		logger.Panic().Err(err)
	}

	// ─── Wait for interrupt (SIGINT/SIGTERM) and shut down gracefully ───────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Warn().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	logger.Warn().Msg("Server exiting")
}
