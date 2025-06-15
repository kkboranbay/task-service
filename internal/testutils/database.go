package testutils

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkboranbay/task-service/internal/config"
	pg "github.com/kkboranbay/task-service/pkg/postgres"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

type TestDB struct {
	Pool   *pgxpool.Pool
	Config config.DatabaseConfig
	dbName string
}

func NewTestDB(t *testing.T) *TestDB {
	t.Helper()

	dbName := fmt.Sprintf("test_%d", time.Now().UnixNano())

	cfg := config.DatabaseConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		User:     getEnvOrDefault("TEST_DB_USER", "postgres"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "postgres"),
		DBName:   "postgres", // пока подключаемся к default database для создания тестовой
		SSLMode:  "disable",
		MaxConns: 5,
		Timeout:  5 * time.Second,
	}

	ctx := context.Background()

	adminPool, err := pg.NewPool(ctx, cfg)
	require.NoError(t, err, "Failed to connect to PostgreSQL")

	_, err = adminPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	require.NoError(t, err, "Failed to create test database")
	adminPool.Close()

	cfg.DBName = dbName
	testPool, err := pg.NewPool(ctx, cfg)
	require.NoError(t, err, "Failed to connect to test database")

	err = runMigrations(ctx, cfg)
	require.NoError(t, err, "Failed to run migrations")

	return &TestDB{
		Pool:   testPool,
		Config: cfg,
		dbName: dbName,
	}
}

func (tdb *TestDB) Close(t *testing.T) {
	t.Helper()

	if tdb.Pool != nil {
		tdb.Pool.Close()
	}

	cfg := tdb.Config
	cfg.DBName = "postgres"

	ctx := context.Background()
	adminPool, err := pg.NewPool(ctx, cfg)
	if err != nil {
		t.Logf("Failed to connect for cleanup: %v", err)
		return
	}
	defer adminPool.Close()

	// Принудительно отключаем все соединения к тестовой БД
	_, err = adminPool.Exec(ctx, fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid()
	`, tdb.dbName))
	if err != nil {
		t.Logf("Failed to terminate connections: %v", err)
	}

	_, err = adminPool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", tdb.dbName))
	if err != nil {
		t.Logf("Failed to drop test database: %v", err)
	}
}

func (tdb *TestDB) Truncate(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	_, err := tdb.Pool.Exec(ctx, "TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
	require.NoError(t, err, "Failed to truncate tables")
}

func runMigrations(ctx context.Context, cfg config.DatabaseConfig) error {
	_, filename, _, _ := runtime.Caller(0) // путь до текущего файла (database.go)
	basePath := filepath.Join(filepath.Dir(filename), "../../migrations")
	sourceURL := "file://" + basePath

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	m, err := migrate.New(sourceURL, connStr)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
