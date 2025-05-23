package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const requestBody = `{
	"title": "Test Newsletter",
	"content": {
		"text": "Newsletter body as plain text",
		"html": "<p>Newsletter body as HTML</p>"
	}
	}`

func createUnconfirmedSubscriber(t *testing.T, app *TestApp) string {
	const body = "name=le%20guin&email=ursula_le_guin%40gmail.com"
	urlChan := make(chan string)
	httpmock.ActivateNonDefault(app.EmailClient.Client())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/email", app.EmailClient.BaseURL()),
		func(r *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			var payload internal.SendEmailRequest
			err = json.Unmarshal(body, &payload)
			assert.Nil(t, err)
			urls := ExtractURLs(payload.HtmlBody)
			require.Equal(t, 1, len(urls))

			urlChan <- urls[0]
			return httpmock.NewStringResponse(http.StatusOK, `{"status": "created"}`), nil
		})
	app.PostSubscriptions(body)
	bodyURL := <-urlChan
	confirmationURL, err := SetURLPort(bodyURL, uint16(app.Port))
	require.Nil(t, err)
	return confirmationURL
}

func createConfirmedSubscriber(t *testing.T, app *TestApp) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	confirmationURL := createUnconfirmedSubscriber(t, app)
	req, err := http.NewRequest("GET", confirmationURL, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "plain/text")
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
}

func TestNewslettersAreNotDeliveredToUncomfirmedSubscribers(t *testing.T) {
	app := SpawnApp()
	_ = createUnconfirmedSubscriber(t, &app)

	httpmock.ActivateNonDefault(app.EmailClient.Client())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/email", app.EmailClient.BaseURL()),
		func(r *http.Request) (*http.Response, error) {
			panic("should not be called")
		})
	resp, err := app.PostNewsletters(requestBody)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestNewslettersAreDeliveredToConfirmedSubscribers(t *testing.T) {
	app := SpawnApp()
	createConfirmedSubscriber(t, &app)
	var reqCnt uint32
	httpmock.ActivateNonDefault(app.EmailClient.Client())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/email", app.EmailClient.BaseURL()),
		func(r *http.Request) (*http.Response, error) {
			atomic.AddUint32(&reqCnt, 1)
			return httpmock.NewStringResponse(http.StatusOK, `{"status": "created"}`), nil
		})
	resp, err := app.PostNewsletters(requestBody)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, uint32(1), atomic.LoadUint32(&reqCnt))
}

func TestNewslettersReturns400ForInvalidData(t *testing.T) {
	app := SpawnApp()
	testCases := []struct {
		body string
		err  string
	}{{
		`{
		"content": {
			"text": "Newsletter body as plain text",
			"html": "<p>Newsletter body as HTML</p>",
		}
		}`,
		"missing title",
	},
		{
			`{"title": "Newsletter!"}`,
			"missing content",
		},
	}
	for _, tc := range testCases {
		resp, err := app.PostNewsletters(tc.body)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, tc.err)
	}
}

func TestRequestsMissingAuthorizationAreRejected(t *testing.T) {
	app := SpawnApp()

	url := fmt.Sprintf("%s/newsletters", app.Address)
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", url, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	authHeader := resp.Header.Get("WWW-Authenticate")
	assert.Equal(t, `Basic realm="publish"`, authHeader)
}

func TestNonExistingUserIsRejected(t *testing.T) {
	app := SpawnApp()

	username := uuid.NewString()
	password := uuid.NewString()
	url := fmt.Sprintf("%s/newsletters", app.Address)
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", url, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	authHeader := resp.Header.Get("WWW-Authenticate")
	assert.Equal(t, `Basic realm="publish"`, authHeader)
}

func TestInvalidPasswordIsRejected(t *testing.T) {
	app := SpawnApp()

	username := app.testUser.Username
	password := uuid.NewString()
	url := fmt.Sprintf("%s/newsletters", app.Address)
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("POST", url, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	authHeader := resp.Header.Get("WWW-Authenticate")
	assert.Equal(t, `Basic realm="publish"`, authHeader)
}
