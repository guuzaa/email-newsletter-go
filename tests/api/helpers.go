package api

import (
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/cmd"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/authentication"
	"github.com/guuzaa/email-newsletter/internal/database"
	"github.com/guuzaa/email-newsletter/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestApp struct {
	Address     string
	Port        uint16
	DBPool      *gorm.DB
	EmailClient *internal.EmailClient
	testUser    *TestUser
	apiClient   *http.Client
}

func (app *TestApp) PostSubscriptions(body string) (*http.Response, error) {
	url := fmt.Sprintf("%s/subscriptions", app.Address)
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return app.apiClient.Do(req)
}

func (app *TestApp) PostNewsletters(body string) (*http.Response, error) {
	url := fmt.Sprintf("%s/newsletters", app.Address)
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(app.testUser.Username, app.testUser.Password)
	return app.apiClient.Do(req)
}

func (app *TestApp) PostLogin(body string) (*http.Response, error) {
	url := fmt.Sprintf("%s/login", app.Address)
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return app.apiClient.Do(req)
}

func (app *TestApp) GetLoginPage() (*http.Response, error) {
	url := fmt.Sprintf("%s/login", app.Address)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return app.apiClient.Do(req)
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
			Host:    "127.0.0.1",
			Port:    0,
			BaseURL: "http://127.0.0.1",
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
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	app := TestApp{
		Address:     fmt.Sprintf("http://%s", settings.Address()),
		Port:        settings.Application.Port,
		DBPool:      nil,
		EmailClient: &emailClient,
		testUser:    GenerateTestUser(),
		apiClient: &http.Client{
			Timeout: 1 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Prevent automatic redirects to handle them manually
				return http.ErrUseLastResponse
			},
			Jar: jar,
		},
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
	srv, err := cmd.Run(settings.Address(), app.DBPool, &emailClient, settings.Application.BaseURL)
	if err != nil {
		panic(err)
	}
	app.Address = fmt.Sprintf("http://%s", srv.Addr)
	if err = app.testUser.Store(app.DBPool); err != nil {
		panic(err)
	}

	u, err := url.Parse(app.Address)
	if err != nil {
		panic(err)
	}
	port, err := strconv.ParseUint(u.Port(), 10, 16)
	if err != nil {
		panic(err)
	}
	app.Port = uint16(port)
	return app
}

func ExtractURLs(text string) []string {
	urlPattern := `(https?://[^\s<>"]+)`
	re := regexp.MustCompile(urlPattern)
	return re.FindAllString(text, -1)
}

func SetURLPort(rawURL string, port uint16) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	newHost := net.JoinHostPort(u.Host, strconv.Itoa(int(port)))
	u.Host = newHost
	return u.String(), nil
}

type TestUser struct {
	UserID   string
	Username string
	Password string
}

func GenerateTestUser() *TestUser {
	return &TestUser{
		UserID:   uuid.NewString(),
		Username: uuid.NewString(),
		Password: uuid.NewString(),
	}
}

func (user *TestUser) Store(db *gorm.DB) error {
	passwordHash, err := authentication.HashPassword(user.Password)
	if err != nil {
		return err
	}

	userModel := &models.User{
		ID:       user.UserID,
		Username: user.Username,
		Password: passwordHash,
	}
	if err := db.Create(userModel).Error; err != nil {
		return err
	}
	return nil
}
