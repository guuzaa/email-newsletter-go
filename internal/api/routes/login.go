package routes

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guuzaa/email-newsletter/internal/api/middleware"
	"github.com/guuzaa/email-newsletter/internal/authentication"
	"github.com/guuzaa/email-newsletter/web"
	"gorm.io/gorm"
)

const (
	flashCookieName = "_flash"
)

type LoginHandler struct {
	db *gorm.DB
}

func NewLoginHandler(db *gorm.DB) *LoginHandler {
	return &LoginHandler{db: db}
}

type FormData struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (h *LoginHandler) get(c *gin.Context) {
	log := middleware.GetContextLogger(c)
	log.Trace().Msg("login page")

	cookie, err := c.Request.Cookie(flashCookieName)
	if err != nil {
		log.Trace().Err(err).Msg("failed to get cookie")
		c.Data(http.StatusOK, "text/html; charset=utf-8", web.LoginHTML)
		return
	}

	c.SetCookie(flashCookieName, "", -1, "/login", "", false, true)
	log.Trace().Str("login error", cookie.Value)
	c.Data(http.StatusOK, "text/html; charset=utf-8", web.LoginHTML)
}

func (h *LoginHandler) post(c *gin.Context) {
	log := middleware.GetContextLogger(c)
	h.db = h.db.WithContext(c.Request.Context())

	var data FormData
	err := c.ShouldBind(&data)
	if err != nil {
		log.Trace().Err(err).Msg("failed to parse request body")
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	crdentials := authentication.Credentials{
		Username: data.Username,
		Password: data.Password,
	}
	if !crdentials.Validate(c, h.db) {
		log.Trace().Msg("failed to validate credentials")
		c.SetCookie(flashCookieName, "invalid credentials", 0, "/login", "", false, true)
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	log.Trace().Msg("login in")
	c.Status(http.StatusOK)
}
