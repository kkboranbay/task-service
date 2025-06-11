package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkboranbay/task-service/internal/api/handler"
	"github.com/kkboranbay/task-service/internal/api/middleware"
	"github.com/kkboranbay/task-service/internal/config"
	"github.com/kkboranbay/task-service/internal/service"
	"github.com/rs/zerolog"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	db         *pgxpool.Pool
	log        *zerolog.Logger
	cfg        config.ServerConfig
}

func NewServer(
	db *pgxpool.Pool,
	taskService *service.TaskService,
	cfg config.Config,
	log *zerolog.Logger,
) *Server {
	router := gin.New()

	jwtMiddleware := middleware.NewJWTMiddleware(cfg.Auth, log)

	healthHandler := handler.NewHealthHandler(db, log)
	healthHandler.Register(router)

	authHandler := handler.NewAuthHandler(jwtMiddleware, log)
	authHandler.Register(router)

	api := router.Group("/api/v1")
	api.Use(jwtMiddleware.AuthRequired())

	taskHandler := handler.NewTaskHandler(taskService, log)
	taskHandler.Register(api)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &Server{
		httpServer: httpServer,
		router:     router,
		db:         db,
		log:        log,
		cfg:        cfg.Server,
	}
}

func (s *Server) Run() error {
	s.log.Info().Str("addr", s.httpServer.Addr).Msg("Запуск HTTP сервера")
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("ошибка запуска HTTP сервера: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info().Msg("Остановка HTTP сервера...")

	shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("ошибка остановки HTTP сервера: %w", err)
	}

	if s.db != nil {
		s.db.Close()
	}

	s.log.Info().Msg("HTTP сервер успешно остановлен")
	return nil
}
