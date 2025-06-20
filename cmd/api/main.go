package main

import (
	"context"
	_ "github.com/kkboranbay/task-service/docs"
	"github.com/kkboranbay/task-service/internal/api"
	"github.com/kkboranbay/task-service/internal/config"
	"github.com/kkboranbay/task-service/internal/repository/postgres"
	"github.com/kkboranbay/task-service/internal/service"
	"github.com/kkboranbay/task-service/pkg/logger"
	pg "github.com/kkboranbay/task-service/pkg/postgres"
	"os"
	"os/signal"
	"syscall"
)

// @title Task Service API
// @version 1.0
// @description REST API для управления задачами с JWT авторизацией
// @termsOfService https://taskservice.com/terms

// @contact.name API Support
// @contact.url https://taskservice.com/support
// @contact.email support@taskservice.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите 'Bearer ' + ваш JWT токен

// @tag.name Authentication
// @tag.description Операции авторизации и аутентификации

// @tag.name Tasks
// @tag.description Операции с задачами (CRUD)

// @tag.name Health
// @tag.description Проверка состояния сервиса

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.L().Fatal().Err(err).Msg("Ошибка загрузки конфигурации")
	}

	logger.SetupLogger(cfg.Logger)
	log := logger.L()
	log.Info().Msg("Запуск сервиса управления задачами")

	ctx := context.Background()
	db, err := pg.NewPool(ctx, cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка подключения к базе данных")
	}
	defer pg.Close(db)

	taskRepo := postgres.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, log)
	server := api.NewServer(db, taskService, *cfg, log)
	go func() {
		if err := server.Run(); err != nil {
			log.Fatal().Err(err).Msg("Ошибка запуска сервера")
		}
	}()
	log.Info().Msgf("Сервис запущен на порту %s", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("Получен сигнал остановки")

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Ошибка при остановке сервера")
	}

	log.Info().Msg("Сервис остановлен")
}
