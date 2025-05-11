package api

import (
	"net/http"
	"testing"

	"github.com/guuzaa/email-newsletter/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestSubscribeReturnsA200forValidFormData(t *testing.T) {
	const body = "name=le%20guin&email=ursula_le_guin%40gmail.com"
	resp, err := app.PostSubscriptions(body)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var subscription models.Subscription
	app.DBPool.First(&subscription)
	assert.Equal(t, "le guin", subscription.Name)
	assert.Equal(t, "ursula_le_guin@gmail.com", subscription.Email)
}

func TestSubscribeReturnsA400WhenDataIsMissing(t *testing.T) {
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
