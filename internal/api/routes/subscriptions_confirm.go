package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/internal/database/models"
	"github.com/guuzaa/email-newsletter/internal/domain"
	"github.com/guuzaa/email-newsletter/internal/middleware"
	"gorm.io/gorm"
)

type ConfirmSubscriptionHandler struct {
	db *gorm.DB
}

func NewConfirmSubscriptionHandler(db *gorm.DB) *ConfirmSubscriptionHandler {
	return &ConfirmSubscriptionHandler{db: db}
}

func (h *ConfirmSubscriptionHandler) confirm(c *gin.Context) {
	log := middleware.GetContextLogger(c)
	h.db.WithContext(c.Request.Context())

	subscriptionToken, ok := c.GetQuery("subscription_token")
	if !ok {
		log.Debug().Msg("missing subscription token")
		c.String(http.StatusBadRequest, "Missing subscription token")
		return
	}

	if !domain.ValidSubscriberToken(subscriptionToken) {
		log.Debug().Msg("invalid subscriber token")
		c.String(http.StatusBadRequest, "Invalid subscription token")
		return
	}

	subscriptionID, err := h.getSubscriberIDFromToken(subscriptionToken)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get subscription ID from token")
		c.String(http.StatusInternalServerError, "Failed to confirm subscription")
		return
	}

	if h.subscriptionHasConfirmed(subscriptionID) {
		log.Trace().Msg("click subscription link twice")
		c.String(http.StatusOK, "You've confirmed the email!")
		return
	}

	if err = h.confirmSubscription(subscriptionID); err != nil {
		log.Debug().Err(err).Msg("failed to confirm subscription")
		c.String(http.StatusInternalServerError, "Failed to confirm subscription")
		return
	}
	log.Trace().Str("subscription ID", subscriptionID).Str("subscription token", subscriptionToken).Msg("subscription confirmed")

	c.String(http.StatusOK, "")
}

func (h *ConfirmSubscriptionHandler) subscriptionHasConfirmed(subscriptionID string) bool {
	var subscription models.Subscription
	result := h.db.Where(&models.Subscription{ID: subscriptionID, Status: models.SubscriptionStatusConfirmed}).First(&subscription)
	return result.Error == nil
}

func (h *ConfirmSubscriptionHandler) confirmSubscription(subscriptionID string) error {
	var subscription models.Subscription
	result := h.db.Model(&subscription).Where("id = ?", subscriptionID).Update("status", models.SubscriptionStatusConfirmed)
	return result.Error
}

func (h *ConfirmSubscriptionHandler) getSubscriberIDFromToken(subscriptionToken string) (string, error) {
	var token models.SubscriptionTokens
	if err := h.db.Where("subscription_token = ?", subscriptionToken).First(&token).Error; err != nil {
		return "", err
	}
	return token.SubscriptionID, nil
}
