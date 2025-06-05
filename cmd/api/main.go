package main

import (
	"context"
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
