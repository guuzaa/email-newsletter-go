package database

import (
	"context"
	"time"

	zerologgorm "github.com/go-mods/zerolog-gorm"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB(settings *internal.Settings) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  settings.PostgresSQLDSN(), // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                      // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{
		Logger: &zerologgorm.GormLogger{
			FieldsExclude: []string{zerologgorm.DurationFieldName, zerologgorm.FileFieldName},
		},
	})
	db = db.WithContext(internal.Logger().WithContext(context.Background()))

	db.AutoMigrate(&models.Subscription{})
	db.AutoMigrate(&models.SubscriptionTokens{})
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return db, err
}
