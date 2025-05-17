package cmd

import (
	"net"
	"net/http"

	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/api/routes"
	"github.com/guuzaa/email-newsletter/internal/database"
	"gorm.io/gorm"
)

var logger = internal.Logger()

func Build(config *internal.Settings) (*http.Server, error) {
	senderEmail, err := config.EmailClient.Sender()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse sender email")
		return nil, err
	}
	timeout := config.EmailClient.Timeout()
	emailClient := internal.NewEmailClient(config.EmailClient.BaseURL, senderEmail, config.EmailClient.AuthorizationToken, timeout)

	db, err := database.SetupDB(config)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect database")
		return nil, err
	}
	return Run(config.Address(), db, &emailClient, config.Application.BaseURL)
}

func Run(address string, db *gorm.DB, emailClient *internal.EmailClient, baseURL string) (*http.Server, error) {
	r := routes.SetupRouter(db, emailClient, baseURL)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create listener")
		return nil, err
	}

	// ─── Start server in its own goroutine ───────────────────────────────────────
	srv := &http.Server{
		Handler: r,
		Addr:    listener.Addr().String(),
	}
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("listen and serve")
		}
	}()
	return srv, nil
}
