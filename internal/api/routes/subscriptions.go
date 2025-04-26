package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/internal/domain"
	"github.com/guuzaa/email-newsletter/internal/middleware"
	"github.com/guuzaa/email-newsletter/internal/models"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	db *gorm.DB
}

func NewSubscriptionHandler(db *gorm.DB) *SubscriptionHandler {
	return &SubscriptionHandler{db: db}
}

func (h *SubscriptionHandler) insertSubscription(c *gin.Context, subscription models.Subscription) error {
	log := middleware.GetContextLogger(c)
	log.Trace().Msg("inserting subscription")
	subscription.ID = uuid.NewString()
	subscription.SubscribedAt = time.Now()

	if err := h.db.Create(&subscription).Error; err != nil {
		log.Warn().Err(err).Msg("failed to create subscription in database")
		return err
	}
	log.Trace().Str("name", subscription.Name).Str("email", subscription.Email).Msg("added new subscriber")
	return nil
}

func (h *SubscriptionHandler) parseSubscription(c *gin.Context) (models.Subscription, error) {
	log := middleware.GetContextLogger(c)

	var data models.Subscription
	if err := c.ShouldBind(&data); err != nil {
		log.Trace().Err(err).Msg("failed to bind request body")
		return models.Subscription{}, err
	}
	name, err := domain.SubscriberNameFrom(data.Name)
	if err != nil {
		log.Trace().Err(err).Msg("failed to parse name")
		return models.Subscription{}, err
	}
	data.Name = name.String()

	email, err := domain.SubscriberEmailFrom(data.Email)
	if err != nil {
		log.Trace().Err(err).Msg("failed to parse email")
		return models.Subscription{}, err
	}
	data.Email = email.String()

	log.Trace().Str("name", data.Name).Str("email", data.Email).Msg("parsed subscription")
	return data, nil
}

func (h *SubscriptionHandler) subscribe(c *gin.Context) {
	log := middleware.GetContextLogger(c)
	subscription, err := h.parseSubscription(c)
	if err != nil {
		log.Trace().Err(err).Msg("failed to parse subscription")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err = h.insertSubscription(c, subscription); err != nil {
		log.Trace().Err(err).Msg("failed to insert subscription")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "")
}
