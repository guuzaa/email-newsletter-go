package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/database/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfirmationsWithoutTokenAreRejectedWithA400(t *testing.T) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	app := SpawnApp()
	url := fmt.Sprintf("%s/subscriptions/confirm", app.Address)
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "plain/text")
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTheLinkReturnedBySubscribeReturnsA200IfCalled(t *testing.T) {
	const body = "name=le%20guin&email=ursula_le_guin%40gmail.com"
	app := SpawnApp()
	httpmock.ActivateNonDefault(app.EmailClient.Client())
	urlChan := make(chan string)
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
			urls := ExtractURLs(payload.HtmlBody)
			require.Equal(t, 1, len(urls))
			assert.Contains(t, urls[0], "127.0.0.1")
			urls = append(urls, ExtractURLs(payload.TextBody)...)
			require.Equal(t, 2, len(urls))
			assert.Equal(t, urls[0], urls[1])

			urlChan <- urls[0]
			return httpmock.NewStringResponse(http.StatusOK, `{"status": "mocked"}`), nil
		})

	go func() {
		resp, err := app.PostSubscriptions(body)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}()

	client := http.Client{
		Timeout: 1 * time.Second,
	}
	bodyURL := <-urlChan
	confirmationURL, err := SetURLPort(bodyURL, uint16(app.Port))
	require.Nil(t, err)
	req, err := http.NewRequest("GET", confirmationURL, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "plain/text")
	confirmResp, err := client.Do(req)
	assert.Nil(t, err)
	defer confirmResp.Body.Close()
	assert.Equal(t, http.StatusOK, confirmResp.StatusCode)
}

func TestClickingOnTheConfirmationLinkConfirmsASubscriber(t *testing.T) {
	const body = "name=le%20guin&email=ursula_le_guin%40gmail.com"
	app := SpawnApp()
	httpmock.ActivateNonDefault(app.EmailClient.Client())
	urlChan := make(chan string)
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
			urls := ExtractURLs(payload.HtmlBody)
			require.Equal(t, 1, len(urls))
			assert.Contains(t, urls[0], "127.0.0.1")
			urls = append(urls, ExtractURLs(payload.TextBody)...)
			require.Equal(t, 2, len(urls))
			assert.Equal(t, urls[0], urls[1])

			urlChan <- urls[0]
			return httpmock.NewStringResponse(http.StatusOK, `{"status": "mocked"}`), nil
		})

	go func() {
		resp, err := app.PostSubscriptions(body)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}()

	client := http.Client{
		Timeout: 1 * time.Second,
	}
	bodyURL := <-urlChan
	confirmationURL, err := SetURLPort(bodyURL, uint16(app.Port))
	require.Nil(t, err)
	req, err := http.NewRequest("GET", confirmationURL, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "plain/text")
	confirmResp, err := client.Do(req)
	assert.Nil(t, err)
	defer confirmResp.Body.Close()
	assert.Equal(t, http.StatusOK, confirmResp.StatusCode)
	var subscription models.Subscription
	app.DBPool.First(&subscription)
	assert.Equal(t, "le guin", subscription.Name)
	assert.Equal(t, "ursula_le_guin@gmail.com", subscription.Email)
	assert.Equal(t, models.SubscriptionStatusConfirmed, subscription.Status)
}

func TestClickingOnTheConfirmationLinkTwiceConfirmsASubscriber(t *testing.T) {
	const body = "name=le%20guin&email=ursula_le_guin%40gmail.com"
	app := SpawnApp()
	httpmock.ActivateNonDefault(app.EmailClient.Client())
	urlChan := make(chan string)
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
			urls := ExtractURLs(payload.HtmlBody)
			require.Equal(t, 1, len(urls))
			assert.Contains(t, urls[0], "127.0.0.1")
			urls = append(urls, ExtractURLs(payload.TextBody)...)
			require.Equal(t, 2, len(urls))
			assert.Equal(t, urls[0], urls[1])

			urlChan <- urls[0]
			return httpmock.NewStringResponse(http.StatusOK, `{"status": "mocked"}`), nil
		})

	go func() {
		resp, err := app.PostSubscriptions(body)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}()

	client := http.Client{
		Timeout: 1 * time.Second,
	}
	bodyURL := <-urlChan
	confirmationURL, err := SetURLPort(bodyURL, uint16(app.Port))
	require.Nil(t, err)
	req, err := http.NewRequest("GET", confirmationURL, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "plain/text")
	confirmResp, err := client.Do(req)
	assert.Nil(t, err)
	defer confirmResp.Body.Close()
	assert.Equal(t, http.StatusOK, confirmResp.StatusCode)

	var subscription models.Subscription
	app.DBPool.First(&subscription)
	assert.Equal(t, "le guin", subscription.Name)
	assert.Equal(t, "ursula_le_guin@gmail.com", subscription.Email)
	assert.Equal(t, models.SubscriptionStatusConfirmed, subscription.Status)

	req, err = http.NewRequest("GET", confirmationURL, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "plain/text")
	confirmResp, err = client.Do(req)
	assert.Nil(t, err)
	defer confirmResp.Body.Close()
	msg, err := io.ReadAll(confirmResp.Body)
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, confirmResp.StatusCode)
	assert.Equal(t, "You've confirmed the email!", string(msg))
}
