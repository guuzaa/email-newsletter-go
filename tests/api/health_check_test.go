package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guuzaa/email-newsletter/internal/api/routes"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	app := SpawnApp()
	r := routes.SetupRouter(app.DBPool)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health_check", nil)
	req.Header.Set("Content-Type", "plain/text")
	r.ServeHTTP(w, req)
	assert.Equal(t, 0, w.Body.Len())
	assert.Equal(t, http.StatusOK, w.Code)
}
