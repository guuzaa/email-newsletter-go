package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/cmd"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestApp struct {
	Address     string
	Port        uint16
	DBPool      *gorm.DB
	EmailClient *internal.EmailClient
}

func (app *TestApp) PostSubscriptions(body string) (*http.Response, error) {
	url := fmt.Sprintf("%s/subscriptions", app.Address)
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return client.Do(req)
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
			Host: "127.0.0.1",
			Port: 0,
		},
		EmailClient: internal.EmailClientSettings{
			BaseURL:             "http://localhost:8081",
			SenderEmail:         "test@example.com",
			AuthorizationToken:  "test_token",
			TimeoutMilliseconds: 1000,
		},
	}

	senderEmail, err := settings.EmailClient.Sender()
	if err != nil {
		panic(err)
	}
	emailClient := internal.NewEmailClient(settings.EmailClient.BaseURL, senderEmail, settings.EmailClient.AuthorizationToken, settings.EmailClient.Timeout())
	app := TestApp{
		Address:     fmt.Sprintf("http://%s", settings.Address()),
		Port:        settings.Application.Port,
		DBPool:      nil,
		EmailClient: &emailClient,
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
	app.DBPool, _ = database.SetupDB(&settings)
	srv, err := cmd.Run(settings.Address(), app.DBPool, &emailClient)
	if err != nil {
		panic(err)
	}
	app.Address = fmt.Sprintf("http://%s", srv.Addr)
	return app
}
