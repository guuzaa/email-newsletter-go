package internal

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guuzaa/email-newsletter/internal/api/routes"
	"github.com/guuzaa/email-newsletter/internal/middleware"
	"github.com/guuzaa/email-newsletter/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var logger = middleware.Logger()

func SetupDB(settings *Settings) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  settings.PostgresSQLDSN(), // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                      // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{}) // TODO: use zerolog as logger

	db.AutoMigrate(&models.Subscription{})
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return db, err
}

func Run(config *Settings) {
	db, err := SetupDB(config)
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
