package routes

import (
	"errors"
	"net/http"

	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/api/middleware"
	"github.com/guuzaa/email-newsletter/internal/database/models"
	"github.com/guuzaa/email-newsletter/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewslettersHandler struct {
	db          *gorm.DB
	emailClient *internal.EmailClient
}

func NewNewslettersHandler(db *gorm.DB, emailClient *internal.EmailClient) *NewslettersHandler {
	return &NewslettersHandler{
		db:          db,
		emailClient: emailClient,
	}
}

type BodyData struct {
	Title   string  `json:"title" binding:"required"`
	Content Content `json:"content" binding:"required"`
}

type Content struct {
	Html string `json:"html" binding:"required"`
	Text string `json:"text" binding:"required"`
}
type Credentials struct {
	Username string
	Password string
}

type ConfirmedSubscriber struct {
	Email domain.SubscriberEmail `gorm:"email"`
}

func (h *NewslettersHandler) basicAuthentication(c *gin.Context) (Credentials, error) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		return Credentials{}, errors.New("missing authorization header")
	}
	return Credentials{
		Username: username,
		Password: password,
	}, nil
}

func (h *NewslettersHandler) publishNewsletter(c *gin.Context) {
	log := middleware.GetContextLogger(c)
	h.db = h.db.WithContext(c.Request.Context())

	_, err := h.basicAuthentication(c)
	if err != nil {
		log.Debug().Err(err).Msg("failed to decode basic auth")
		c.Header("WWW-Authenticate", `Basic realm="publish"`)
		c.String(http.StatusUnauthorized, "Invalid credentials")
		return
	}

	var body BodyData
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Trace().Err(err).Msg("failed to bind request body")
		c.String(http.StatusBadRequest, "")
		return
	}
	confirmedSubscribers := h.getConfirmedSubscribers()
	log.Debug().Int("len confirmed subscribers", len(confirmedSubscribers)).Send()
	for _, subscriber := range confirmedSubscribers {
		if err := h.emailClient.SendEmail(subscriber.Email, body.Title, body.Content.Html, body.Content.Text); err != nil {
			log.Warn().Err(err).Str("email", subscriber.Email.String()).Msg("failed to send email")
			c.String(http.StatusInternalServerError, "Failed to send email")
			return
		}
		log.Trace().Msgf("sending email to %s", subscriber.Email)
	}
	c.String(http.StatusOK, "")
}

func (h *NewslettersHandler) getConfirmedSubscribers() []ConfirmedSubscriber {
	var confirmedSubscribers []ConfirmedSubscriber
	var subscriberEmails []string
	if err := h.db.Model(&models.Subscription{}).Select("email").Where("status = ?", models.SubscriptionStatusConfirmed).Find(&subscriberEmails).Error; err != nil {
		return confirmedSubscribers
	}
	for _, subscriberEmail := range subscriberEmails {
		email, err := domain.SubscriberEmailFrom(subscriberEmail)
		if err != nil {
			continue
		}
		confirmedSubscribers = append(confirmedSubscribers, ConfirmedSubscriber{Email: email})
	}
	return confirmedSubscribers
}
