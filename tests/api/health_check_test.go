package api

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	app := SpawnApp()
	url := fmt.Sprintf("%s/health_check", app.Address)
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "plain/text")
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(body))
}
