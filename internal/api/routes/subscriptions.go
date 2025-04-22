package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *SubscriptionHandler) subscribe(c *gin.Context) {
	r := middleware.Logger()
	requestID := c.Value("requestID")
	var data models.Subscription
	if err := c.ShouldBind(&data); err != nil {
		r.Trace().Err(err).Msgf("requestID: %s, failed to bind request body", requestID)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	switch {
	case len(data.Name) == 0 && len(data.Email) == 0:
		r.Trace().Msgf("requestID: %s, missing both name and email", requestID)
		c.String(http.StatusBadRequest, "missing both name and email")
	case len(data.Name) == 0:
		r.Trace().Msgf("requestID: %s, missing the name", requestID)
		c.String(http.StatusBadRequest, "missing the name")
	case len(data.Email) == 0:
		r.Trace().Msgf("requestID: %s, missing the email", requestID)
		c.String(http.StatusBadRequest, "missing the email")
	default:
		r.Trace().Msgf("requestID: %s, creating subscription", requestID)
		data.ID = uuid.NewString()
		data.SubscribedAt = time.Now()
		if err := h.db.Create(&data).Error; err != nil {
			r.Warn().Err(err).Msgf("requestID: %s, failed to create subscription in database", requestID)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		r.Trace().Msgf("requestID: %s, Adding '%s' '%s' as a new subscriber", requestID, data.Name, data.Email)
		c.String(http.StatusOK, "")
	}
}
