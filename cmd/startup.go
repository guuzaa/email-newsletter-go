package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/api/routes"
	"github.com/guuzaa/email-newsletter/internal/database"
)

var logger = internal.Logger()

func Run(config *internal.Settings) {
	senderEmail, err := config.EmailClient.Sender()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse sender email")
	}
	timeout := config.EmailClient.Timeout()
	_ = internal.NewEmailClient(config.EmailClient.BaseURL, senderEmail, config.EmailClient.AuthorizationToken, timeout)

	db, err := database.SetupDB(config)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect database")
	}

	r := routes.SetupRouter(db)

	// ─── Start server in its own goroutine ───────────────────────────────────────
	srv := &http.Server{
		Addr:    config.Address(),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("listen and serve")
		}
	}()

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
