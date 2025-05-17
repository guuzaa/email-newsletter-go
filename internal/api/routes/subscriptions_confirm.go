package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/internal/database/models"
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
	subscriptionToken, ok := c.GetQuery("subscription_token")
	if !ok {
		log.Debug().Msg("Missing subscription token")
		c.String(http.StatusBadRequest, "Missing subscription token")
		return
	}

	subscriptionID, err := h.getSubscriberIDFromToken(subscriptionToken)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get subscription ID from token")
		c.String(http.StatusInternalServerError, "Failed to confirm subscription")
		return
	}

	if err = h.confirmSubscription(subscriptionID); err != nil {
		log.Debug().Err(err).Msg("Failed to confirm subscription")
		c.String(http.StatusInternalServerError, "Failed to confirm subscription")
		return
	}

	c.String(http.StatusOK, "")
}

func (h *ConfirmSubscriptionHandler) confirmSubscription(subscriptionID string) error {
	var subscription models.Subscription
	result := h.db.Model(&subscription).Where("id = ?", subscriptionID).Update("status", "confirmed")
	return result.Error
}

func (h *ConfirmSubscriptionHandler) getSubscriberIDFromToken(subscriptionToken string) (string, error) {
	var token models.SubscriptionTokens
	if err := h.db.Where("subscription_token = ?", subscriptionToken).First(&token).Error; err != nil {
		return "", err
	}
	return token.SubscriptionID, nil
}
