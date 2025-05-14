package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/guuzaa/email-newsletter/internal/domain"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func subject() string {
	return faker.New().Lorem().Sentence(2)
}

func content() string {
	return faker.New().Lorem().Sentence(10)
}

func email() domain.SubscriberEmail {
	subscriberEmail, err := domain.SubscriberEmailFrom(faker.New().Internet().Email())
	if err != nil {
		panic(err)
	}
	return subscriberEmail
}

func emailClient(baseUrl string) *EmailClient {
	token := faker.New().RandomStringWithLength(5)
	client := NewEmailClient(baseUrl, email(), token, 200*time.Millisecond)
	return &client
}

func TestSendEmailSendsTheExpectedRequest(t *testing.T) {
	reqCnt := uint32(0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint32(&reqCnt, 1)

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if token := r.Header.Get("X-Postmark-Server-Token"); token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			http.Error(w, "Invalid content type", http.StatusUnsupportedMediaType)
			return
		}
		if !strings.Contains(r.URL.String(), "/email") {
			http.Error(w, "Invalid url", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		var payload SendEmailRequest
		err = json.Unmarshal(body, &payload)
		if err != nil {
			http.Error(w, "Failed to unmarshal request body", http.StatusInternalServerError)
			return
		}
		assert.NotEmpty(t, payload.From)
		assert.NotEmpty(t, payload.To)
		assert.NotEmpty(t, payload.Subject)
		assert.NotEmpty(t, payload.HtmlBody)
		assert.NotEmpty(t, payload.TextBody)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "mock response"}`))
	}))
	defer server.Close()

	emailClient := emailClient(server.URL)
	subscriberEmail := email()
	content := content()
	err := emailClient.SendEmail(subscriberEmail, subject(), content, content)
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), reqCnt)
}

func TestSendEmailFailsIfTheServerReturns500(t *testing.T) {
	reqCnt := uint32(0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint32(&reqCnt, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	emailClient := emailClient(server.URL)
	subscriberEmail := email()
	subject := subject()
	content := content()
	err := emailClient.SendEmail(subscriberEmail, subject, content, content)
	assert.NotNil(t, err)
	assert.Equal(t, uint32(1), reqCnt)
}

func TestSendEmailTimesoutIfTheServerTakesTooLong(t *testing.T) {
	reqCnt := uint32(0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint32(&reqCnt, 1)
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "mock response"}`))
	}))
	defer server.Close()

	emailClient := emailClient(server.URL)
	subscriberEmail := email()
	subject := subject()
	content := content()
	err := emailClient.SendEmail(subscriberEmail, subject, content, content)
	assert.NotNil(t, err)
	assert.Equal(t, uint32(1), reqCnt)
}
