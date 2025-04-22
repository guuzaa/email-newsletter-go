package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guuzaa/email-newsletter/internal/logger"
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
	r := logger.Logger()
	var data models.Subscription
	if err := c.ShouldBind(&data); err != nil {
		r.Trace().Err(err).Msg("failed to bind request body")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	switch {
	case len(data.Name) == 0 && len(data.Email) == 0:
		r.Trace().Msg("missing both name and email")
		c.String(http.StatusBadRequest, "missing both name and email")
	case len(data.Name) == 0:
		r.Trace().Msg("missing the name")
		c.String(http.StatusBadRequest, "missing the name")
	case len(data.Email) == 0:
		r.Trace().Msg("missing the email")
		c.String(http.StatusBadRequest, "missing the email")
	default:
		r.Trace().Msg("creating subscription")
		data.ID = uuid.NewString()
		data.SubscribedAt = time.Now()
		if err := h.db.Create(&data).Error; err != nil {
			r.Warn().Err(err).Msg("failed to create subscription in database")
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		r.Trace().Msg("subscription created successfully")
		c.String(http.StatusOK, "")
	}
}
