package main

import (
	"context"
	"log"
	"time"

	"github.com/Ainyx-backend/config"
	"github.com/Ainyx-backend/db/sqlc"
	"github.com/Ainyx-backend/internal/handler"
	"github.com/Ainyx-backend/internal/logger"
	"github.com/Ainyx-backend/internal/middleware"
	"github.com/Ainyx-backend/internal/repository"
	"github.com/Ainyx-backend/internal/routes"
	"github.com/Ainyx-backend/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	zLogger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer zLogger.Sync()

	if cfg.DBURL == "" {
		zLogger.Fatal("DB_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		zLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		zLogger.Fatal("failed to ping database", zap.Error(err))
	}

	queries := sqlc.New(pool)
	repo := repository.NewUserRepository(queries)
	userService := service.NewUserService(repo)
	validate := validator.New()
	userHandler := handler.NewUserHandler(userService, validate, zLogger)

	app := fiber.New(fiber.Config{
		AppName: "Ainyx Backend",
	})

	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(zLogger))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	routes.Setup(app, userHandler)

	zLogger.Info("server starting", zap.String("port", cfg.Port))
	if err := app.Listen(":" + cfg.Port); err != nil {
		zLogger.Fatal("server stopped", zap.Error(err))
	}
}