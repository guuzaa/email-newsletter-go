package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/models"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeReturnsA200forValidFormData(t *testing.T) {
	const body = "name=le%20guin&email=ursula_le_guin%40gmail.com"
	app := SpawnApp()
	httpmock.ActivateNonDefault(app.EmailClient.Client())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/email", app.EmailClient.BaseURL()),
		func(r *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			var payload internal.SendEmailRequest
			err = json.Unmarshal(body, &payload)
			assert.Nil(t, err)
			assert.NotEmpty(t, payload.From)
			assert.NotEmpty(t, payload.To)
			assert.Equal(t, "Welcome!", payload.Subject)
			assert.NotEmpty(t, payload.TextBody)
			assert.NotEmpty(t, payload.HtmlBody)
			return httpmock.NewStringResponse(http.StatusOK, `{"status": "created"}`), nil
		})

	resp, err := app.PostSubscriptions(body)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var subscription models.Subscription
	app.DBPool.First(&subscription)
	assert.Equal(t, "le guin", subscription.Name)
	assert.Equal(t, "ursula_le_guin@gmail.com", subscription.Email)
}

func TestSubscribePersistsTheNewSubscriber(t *testing.T) {
	const body = "name=le%20guin&email=ursula_le_guin%40gmail.com"
	app := SpawnApp()
	httpmock.ActivateNonDefault(app.EmailClient.Client())
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/email", app.EmailClient.BaseURL()),
		func(r *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			var payload internal.SendEmailRequest
			err = json.Unmarshal(body, &payload)
			assert.Nil(t, err)
			assert.NotEmpty(t, payload.From)
			assert.NotEmpty(t, payload.To)
			assert.Equal(t, "Welcome!", payload.Subject)
			assert.NotEmpty(t, payload.TextBody)
			assert.NotEmpty(t, payload.HtmlBody)
			return httpmock.NewStringResponse(http.StatusOK, `{"status": "created"}`), nil
		})

	resp, err := app.PostSubscriptions(body)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var subscription models.Subscription
	app.DBPool.First(&subscription)
	assert.Equal(t, "le guin", subscription.Name)
	assert.Equal(t, "ursula_le_guin@gmail.com", subscription.Email)
	assert.Equal(t, "pending_confirmation", subscription.Status)
}

func TestSubscribeReturnsA400WhenDataIsMissing(t *testing.T) {
	app := SpawnApp()
	var testCases = []struct {
		name     string
		body     string
		expected int
	}{
		{"missing the name", "email=ursula_le_guin%40gmail.com", http.StatusBadRequest},
		{"missing the email", "name=le%20guin", http.StatusBadRequest},
		{"missing both name and email", "", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		resp, err := app.PostSubscriptions(tc.body)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, tc.expected, resp.StatusCode)
	}
}

func TestSubscribeReturnsA200WhenFieldsArePresentButEmpty(t *testing.T) {
	app := SpawnApp()
	var testCases = []struct {
		name     string
		body     string
		expected int
	}{
		{"empty name", "name=&email=ursula_le_guin%40gmail.com", http.StatusBadRequest},
		{"empty email", "name=le%20guin&email=", http.StatusBadRequest},
		{"invalid email", "name=le%20guin&email=invalid-email", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		resp, err := app.PostSubscriptions(tc.body)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, tc.expected, resp.StatusCode)
	}
}
