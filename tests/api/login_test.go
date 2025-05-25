package api

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnErrorFlashMessageIsSetOnFailure(t *testing.T) {
	app := SpawnApp()
	loginBody := `
	{
		"username": "random-username",
		"password": "random-password"
	}`
	resp, err := app.PostLogin(loginBody)
	require.Nil(t, err)
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
	assert.Equal(t, "/login", resp.Header.Get("Location"))
	defer resp.Body.Close()
	cookies := resp.Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "_flash", cookies[0].Name)
	assert.Equal(t, "invalid+credentials", cookies[0].Value)

	resp, err = app.GetLoginPage()
	require.Nil(t, err)
	require.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err)
	htmlPage := string(body)
	assert.Nil(t, err)
	assert.NotEmpty(t, htmlPage)
	// assert.Contains(t, htmlPage, `<p><i>invalid crdentials</i></p>`)
}
