package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"net/http"
)

type HealthHandler struct {
	db  *pgxpool.Pool
	log *zerolog.Logger
}

func NewHealthHandler(db *pgxpool.Pool, log *zerolog.Logger) *HealthHandler {
	return &HealthHandler{
		db:  db,
		log: log,
	}
}

func (h *HealthHandler) Register(router *gin.Engine) {
	router.GET("/health", h.Check)
	router.GET("/readiness", h.Readiness)
}

type HealthResponse struct {
	Status string `json:"status"`
}

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}

// Readiness проверяет готовность сервиса к работе, включая соединение с БД
func (h *HealthHandler) Readiness(c *gin.Context) {
	if err := h.db.Ping(c.Request.Context()); err != nil {
		h.log.Error().Err(err).Msg("ошибка подключения к базе данных")
		c.JSON(http.StatusServiceUnavailable, HealthResponse{Status: "database connection failed"})
		return
	}

	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
