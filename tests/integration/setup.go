package integration

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestApp struct {
	Address string
	DBPool  *gorm.DB
}

func SpawnApp() TestApp {
	settings := internal.Settings{
		Database: internal.DatabaseSettings{
			Host:         "localhost",
			Port:         5432,
			DatabaseName: "postgres",
			Username:     "postgres",
			Password:     "password",
		},
		Application: internal.ApplicationSettings{
			Host: "localhost",
			Port: 8080,
		},
	}
	app := TestApp{
		Address: settings.Address(),
		DBPool:  nil,
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  settings.PostgresSQLDSN(), // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                      // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	settings.Database.DatabaseName = uuid.NewString()
	createQuery := fmt.Sprintf(`CREATE DATABASE "%s"`, settings.Database.DatabaseName)
	if result := db.Exec(createQuery); result.Error != nil {
		panic(result.Error)
	}
	app.DBPool, _ = internal.SetupDB(&settings)
	return app
}
