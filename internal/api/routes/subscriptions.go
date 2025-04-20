package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	var data models.Subscription
	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	switch {
	case len(data.Name) == 0 && len(data.Email) == 0:
		c.String(http.StatusBadRequest, "missing both name and email")
	case len(data.Name) == 0:
		c.String(http.StatusBadRequest, "missing the name")
	case len(data.Email) == 0:
		c.String(http.StatusBadRequest, "missing the email")
	default:
		data.ID = uuid.NewString()
		data.SubscribedAt = time.Now()
		if err := h.db.Create(&data).Error; err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusOK, "")
	}
}
