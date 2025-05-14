package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/guuzaa/email-newsletter/internal/domain"
)

type EmailClient struct {
	httpClient         *http.Client
	baseUrl            string
	sender             domain.SubscriberEmail
	authorizationToken string
}

func NewEmailClient(baseUrl string, sender domain.SubscriberEmail, authorizationToken string, timeout time.Duration) EmailClient {
	return EmailClient{
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:    100, // Connection pool size
				IdleConnTimeout: 30 * time.Second,
			},
		},
		baseUrl:            baseUrl,
		sender:             sender,
		authorizationToken: authorizationToken,
	}
}

type SendEmailRequest struct {
	From     string `json:"From"`
	To       string `json:"To"`
	Subject  string `json:"Subject"`
	HtmlBody string `json:"HtmlBody"`
	TextBody string `json:"TextBody"`
}

func (ec *EmailClient) SendEmail(recipient domain.SubscriberEmail, subject, htmlContent, textContent string) error {
	url := fmt.Sprintf("%s/email", ec.baseUrl)
	request := SendEmailRequest{
		From:     ec.sender.String(),
		To:       recipient.String(),
		Subject:  subject,
		HtmlBody: htmlContent,
		TextBody: textContent,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("X-Postmark-Server-Token", ec.authorizationToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ec.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (ec *EmailClient) Client() *http.Client {
	return ec.httpClient
}

func (ec *EmailClient) BaseURL() string {
	return ec.baseUrl
}
