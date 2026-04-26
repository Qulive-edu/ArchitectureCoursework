// backend/cmd/app/main.go
package main

import (
	"backend/config"
	"backend/internal/app"
	"backend/internal/entity"
	"backend/internal/infrastructure/postgres"
	"backend/internal/infrastructure/redis"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/crypto/bcrypt"
)

const (
	ModeServer      = "server"
	ModeMigrate     = "migrate"
	ModeCreateAdmin = "create-admin"
	ModeClearCache  = "clear-cache"
)

func main() {
	mode := flag.String("mode", ModeServer, "Режим запуска: server|migrate|create-admin|clear-cache")
	configPath := flag.String("config", "", "Путь к config.yml")

	adminEmail := flag.String("admin-email", "", "Email администратора")
	adminName := flag.String("admin-name", "Admin", "Имя администратора")
	adminPassword := flag.String("admin-password", "", "Пароль администратора")

	migrateUp := flag.Bool("migrate-up", true, "Применить миграции вверх")
	migrateSteps := flag.Int("migrate-steps", 0, "Количество шагов миграции (0 = все)")

	cachePattern := flag.String("cache-pattern", "*", "Pattern для очистки кэша")

	flag.Parse()

	if *configPath == "" {
		*configPath = os.Getenv("CONFIG_PATH")
	}
	if *configPath == "" {
		*configPath = "config/config.yml"
	}

	cfg := config.NewConfigWithFile(*configPath)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	switch *mode {
	case ModeServer:
		runServer(context.Background(), cfg, logger)
	case ModeMigrate:
		runMigrations(cfg, logger, *migrateUp, *migrateSteps)
	case ModeCreateAdmin:
		createAdminUser(cfg, logger, *adminEmail, *adminName, *adminPassword)
	case ModeClearCache:
		clearCache(cfg, logger, *cachePattern)
	default:
		logger.Error("unknown mode", "mode", *mode)
		os.Exit(1)
	}
}

func runServer(ctx context.Context, cfg config.Config, logger *slog.Logger) {
	// Инициализация Redis
	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		logger.Error("redis connect: " + err.Error())
		os.Exit(1)
	}
	defer rdb.Close()

	// Запуск основного приложения
	app.Run(ctx, cfg, logger, rdb)
}

func runMigrations(cfg config.Config, logger *slog.Logger, up bool, steps int) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
	)

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		logger.Error("init migrate: " + err.Error())
		os.Exit(1)
	}
	defer m.Close()

	if up {
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			logger.Error("migration up: " + err.Error())
			os.Exit(1)
		}
		logger.Info("migrations applied successfully")
	} else {
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			logger.Error("migration down: " + err.Error())
			os.Exit(1)
		}
		logger.Info("migrations reverted successfully")
	}
}

func createAdminUser(cfg config.Config, logger *slog.Logger, email, name, password string) {
	if email == "" || password == "" {
		logger.Error("admin-email and admin-password are required")
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := postgres.New(ctx, cfg.Database, logger)
	if err != nil {
		logger.Error("db connect: " + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepo(db)

	// Проверка существования
	_, err = userRepo.GetByEmail(ctx, email)
	if err == nil {
		logger.Warn("user already exists", "email", email)
		os.Exit(0)
	}

	// Создание пользователя
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("hash password: " + err.Error())
		os.Exit(1)
	}

	user := &entity.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := userRepo.Create(ctx, user); err != nil {
		logger.Error("create admin: " + err.Error())
		os.Exit(1)
	}

	logger.Info("admin user created", "email", email, "id", user.ID)
}

func clearCache(cfg config.Config, logger *slog.Logger, pattern string) {
	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		logger.Error("redis connect: " + err.Error())
		os.Exit(1)
	}
	defer rdb.Close()

	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	count := 0
	for iter.Next(ctx) {
		if err := rdb.Del(ctx, iter.Val()).Err(); err != nil {
			logger.Warn("delete key", "key", iter.Val(), "error", err)
			continue
		}
		count++
	}
	if err := iter.Err(); err != nil {
		logger.Error("scan error: " + err.Error())
		os.Exit(1)
	}
	logger.Info("cache cleared", "keys_deleted", count)
}
