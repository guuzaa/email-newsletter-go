package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/domain"
	"github.com/guuzaa/email-newsletter/internal/middleware"
	"github.com/guuzaa/email-newsletter/internal/models"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	db          *gorm.DB
	emailClient *internal.EmailClient
}

func NewSubscriptionHandler(db *gorm.DB, emailClient *internal.EmailClient) *SubscriptionHandler {
	return &SubscriptionHandler{db: db, emailClient: emailClient}
}

func (h *SubscriptionHandler) insertSubscriber(c *gin.Context, subscriber domain.NewSubscriber) error {
	log := middleware.GetContextLogger(c)
	log.Trace().Msg("inserting subscription")
	subscription := models.Subscription{
		Name:   subscriber.Name.String(),
		Email:  subscriber.Email.String(),
		ID:     uuid.NewString(),
		Status: "pending_confirmation",
	}

	if err := h.db.Create(&subscription).Error; err != nil {
		log.Warn().Err(err).Msg("failed to create subscription in database")
		return err
	}
	log.Trace().Str("name", subscription.Name).Str("email", subscription.Email).Msg("added new subscriber")
	return nil
}

func (h *SubscriptionHandler) parseSubscription(c *gin.Context) (domain.NewSubscriber, error) {
	log := middleware.GetContextLogger(c)

	var data models.Subscription
	if err := c.ShouldBind(&data); err != nil {
		log.Trace().Err(err).Msg("failed to bind request body")
		return domain.NewSubscriber{}, err
	}
	name, err := domain.SubscriberNameFrom(data.Name)
	if err != nil {
		log.Trace().Err(err).Msg("failed to parse name")
		return domain.NewSubscriber{}, err
	}

	email, err := domain.SubscriberEmailFrom(data.Email)
	if err != nil {
		log.Trace().Err(err).Msg("failed to parse email")
		return domain.NewSubscriber{}, err
	}

	log.Trace().Str("name", data.Name).Str("email", data.Email).Msg("parsed subscription")
	return domain.NewSubscriber{
		Name:  name,
		Email: email,
	}, nil
}

func (h *SubscriptionHandler) subscribe(c *gin.Context) {
	log := middleware.GetContextLogger(c)
	newSubscriber, err := h.parseSubscription(c)
	if err != nil {
		log.Trace().Err(err).Msg("failed to parse subscription")
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err = h.insertSubscriber(c, newSubscriber); err != nil {
		log.Debug().Err(err).Msg("failed to insert subscription")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err = h.sendConfirmationEmail(newSubscriber); err != nil {
		log.Warn().Err(err).Msg("failed to send confirmation email")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "")
}

func (h *SubscriptionHandler) sendConfirmationEmail(newSubscriber domain.NewSubscriber) error {
	subject := "Welcome!"
	confirmationLink := "https://there-is-no-such-domain.com/subscriptions/confirm"
	htmlContent := fmt.Sprintf(`Welcome to our newsletter!<br />
	Click <a href="%s">here</a> to confirm your subscription.`, confirmationLink)
	textContent := fmt.Sprintf(`Welcome to our newsletter!
	Click %s to confirm your subscription.`, confirmationLink)
	return h.emailClient.SendEmail(newSubscriber.Email, subject, htmlContent, textContent)
}
