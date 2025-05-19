package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/internal"
	"github.com/guuzaa/email-newsletter/internal/api/middleware"
	"github.com/guuzaa/email-newsletter/internal/database/models"
	"github.com/guuzaa/email-newsletter/internal/domain"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	db          *gorm.DB
	emailClient *internal.EmailClient
	baseURL     string
}

func NewSubscriptionHandler(db *gorm.DB, emailClient *internal.EmailClient, baseURL string) *SubscriptionHandler {
	return &SubscriptionHandler{db: db, emailClient: emailClient, baseURL: baseURL}
}

func (h *SubscriptionHandler) insertSubscriber(c *gin.Context, tx *gorm.DB, subscriber domain.NewSubscriber) (string, error) {
	log := middleware.GetContextLogger(c)
	log.Trace().Msg("inserting subscription")
	subscriberID := uuid.NewString()
	subscription := models.Subscription{
		Name:   subscriber.Name.String(),
		Email:  subscriber.Email.String(),
		ID:     subscriberID,
		Status: models.SubscriptionStatusPending,
	}

	if err := tx.Create(&subscription).Error; err != nil {
		return "", err
	}
	log.Trace().Str("name", subscription.Name).Str("email", subscription.Email).Msg("added new subscriber")
	return subscriberID, nil
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

func (h *SubscriptionHandler) hasPendingSubscriber(subscriber domain.NewSubscriber) bool {
	email := subscriber.Email.String()
	var subscription models.Subscription
	result := h.db.Where("email = ?", email).First(&subscription)
	if err := result.Error; err != nil {
		return false
	}
	return result.RowsAffected == 1 && subscription.Status == models.SubscriptionStatusPending
}

func (h *SubscriptionHandler) subscribe(c *gin.Context) {
	log := middleware.GetContextLogger(c)
	h.db = h.db.WithContext(c.Request.Context())

	newSubscriber, err := h.parseSubscription(c)
	if err != nil {
		log.Trace().Err(err).Msg("failed to parse subscription")
		c.String(http.StatusBadRequest, "Invalid subscription")
		return
	}

	if h.hasPendingSubscriber(newSubscriber) {
		log.Trace().Msg("subscribe twice")
		c.String(http.StatusOK, "You've subscribed already!")
		return
	}

	tx := h.db.Begin()
	subscriberID, err := h.insertSubscriber(c, tx, newSubscriber)
	if err != nil {
		log.Debug().Err(err).Msg("failed to insert subscription")
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Failed to store subscription")
		return
	}

	subscriptionToken := domain.NewSubscriptionToken()
	if err = h.storeToken(tx, subscriberID, subscriptionToken); err != nil {
		log.Warn().Err(err).Msg("failed to store subscription")
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Failed to store subscription token")
		return
	}
	log.Debug().Msgf("subscription token stored for subscriber ID %s, token %s", subscriberID, subscriptionToken)

	if err = tx.Commit().Error; err != nil {
		log.Warn().Err(err).Msg("failed to commit transaction")
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Internal server error from database")
		return
	}
	log.Debug().Msgf("subscription created, ID %s, token %s", subscriberID, subscriptionToken)

	if err = h.sendConfirmationEmail(newSubscriber, subscriptionToken); err != nil {
		log.Warn().Err(err).Msg("failed to send confirmation email")
		c.String(http.StatusInternalServerError, "Failed to send confirmation email")
		return
	}
	log.Debug().Msgf("confirmation email sent to subscriber ID %s, token %s", subscriberID, subscriptionToken)

	c.String(http.StatusOK, "")
}

func (h *SubscriptionHandler) sendConfirmationEmail(newSubscriber domain.NewSubscriber, token string) error {
	subject := "Welcome!"
	confirmationLink := fmt.Sprintf("%s/subscriptions/confirm?subscription_token=%s", h.baseURL, token)
	htmlContent := fmt.Sprintf(`Welcome to our newsletter!<br />
	Click <a href="%s">here</a> to confirm your subscription.`, confirmationLink)
	textContent := fmt.Sprintf(`Welcome to our newsletter!
	Click %s to confirm your subscription.`, confirmationLink)
	return h.emailClient.SendEmail(newSubscriber.Email, subject, htmlContent, textContent)
}

func (h *SubscriptionHandler) storeToken(tx *gorm.DB, subscriberID string, subscriptionToken string) error {
	token := models.SubscriptionTokens{
		SubscriptionID:    subscriberID,
		SubscriptionToken: subscriptionToken,
	}
	result := tx.Create(token)
	return result.Error
}
