package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kkboranbay/task-service/internal/api/middleware"
	"github.com/kkboranbay/task-service/internal/model"
	"github.com/rs/zerolog"
	"net/http"
)

type AuthHandler struct {
	jwtMiddleware *middleware.JWTMiddleware
	log           *zerolog.Logger
}

func NewAuthHandler(jwtMiddleware *middleware.JWTMiddleware, log *zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		jwtMiddleware: jwtMiddleware,
		log:           log,
	}
}

func (h *AuthHandler) Register(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("ошибка разбора JSON")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "некорректные данные запроса",
		})
		return
	}

	if req.Username != "admin" || req.Password != "admin" {
		h.log.Warn().Str("username", req.Username).Msg("неверные учетные данные")
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "неверные учетные данные",
		})
		return
	}

	userID := int64(1)
	token, err := h.jwtMiddleware.GenerateToken(userID)
	if err != nil {
		h.log.Error().Err(err).Int64("user_id", userID).Msg("ошибка генерации токена")
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "не удалось сгенерировать токен",
		})
		return
	}

	h.log.Info().Str("username", req.Username).Int64("user_id", userID).Msg("успешная авторизация")

	c.JSON(http.StatusOK, model.LoginResponse{
		Token:  token,
		UserID: userID,
	})
}
