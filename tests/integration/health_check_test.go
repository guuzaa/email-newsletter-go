package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/guuzaa/email-newsletter/internal/api/routes"
	"github.com/guuzaa/email-newsletter/internal/models"

	"github.com/stretchr/testify/assert"
)

var (
	app = SpawnApp()
	r   = routes.SetupRouter(app.DBPool)
)

func TestHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health_check", nil)
	req.Header.Set("Content-Type", "plain/text")
	r.ServeHTTP(w, req)
	assert.Equal(t, 0, w.Body.Len())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSubscribeReturnsA200forValidFormData(t *testing.T) {
	const body = "name=le%20guin&email=ursula_le_guin%40gmail.com"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/subscriptions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var subscription models.Subscription
	app.DBPool.First(&subscription)
	assert.Equal(t, "le guin", subscription.Name)
	assert.Equal(t, "ursula_le_guin@gmail.com", subscription.Email)
}

func TestSubscribeReturnsA400WhenDataIsMissing(t *testing.T) {
	var testCases = []struct {
		errorMsg string
		body     string
		expected int
	}{
		{"missing the name", "email=ursula_le_guin%40gmail.com", http.StatusBadRequest},
		{"missing the email", "name=le%20guin", http.StatusBadRequest},
		{"missing both name and email", "", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/subscriptions", strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(w, req)
		assert.Equal(t, tc.expected, w.Code)
		assert.Equal(t, tc.errorMsg, w.Body.String(), "The API did not fail with 400 Bad Request when the payload was %s", w.Body.String())
	}
}
