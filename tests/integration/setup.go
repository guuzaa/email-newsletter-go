package integration

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/api/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestApp struct {
	Address string
	DBPool  *gorm.DB
}

func SpawnApp() TestApp {
	settings := internal.Settings{
		DatabaseSettings: internal.DatabaseSettings{
			Host:         "localhost",
			Port:         5432,
			DatabaseName: "postgres",
			Username:     "postgres",
			Password:     "password",
		},
		ApplicationPort: 8080,
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

	settings.DatabaseName = uuid.NewString()
	createQuery := fmt.Sprintf(`CREATE DATABASE "%s"`, settings.DatabaseName)
	if result := db.Exec(createQuery); result.Error != nil {
		panic(result.Error)
	}
	app.DBPool, _ = internal.SetupDB(&settings)
	return app
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	app := SpawnApp()
	return routes.SetupRouter(app.DBPool)
}
