package internal

import (
	"time"

	"github.com/guuzaa/email-newsletter/internal/api/routes"
	"github.com/guuzaa/email-newsletter/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB(settings *Settings) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  settings.PostgresSQLDSN(), // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                      // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})

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
		panic("failed to connect database")
	}

	r := routes.SetupRouter(db)
	r.Run(config.Address())
}
