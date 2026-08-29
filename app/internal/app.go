package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/adriein/hastypal/database"
	"github.com/adriein/hastypal/internal/telegram"
	"github.com/adriein/hastypal/internal/whatsapp"
	"github.com/adriein/hastypal/pkg/constants"
	"github.com/adriein/hastypal/pkg/helper"
	"github.com/adriein/hastypal/pkg/logger"
	"github.com/joho/godotenv"
)

type ShutdownFunc func(context.Context) error

type Modules struct {
	Database *sql.DB
	Logger   *slog.Logger
	Telegram telegram.TelegramService
	Whatsapp whatsapp.WhatsappService
}

type App struct {
	Modules  *Modules
	Shutdown ShutdownFunc
}

func NewApp() *App {
	if os.Getenv(constants.Env) != constants.Pro {
		dotenvErr := godotenv.Load()

		if dotenvErr != nil {
			log.Fatal("Error loading .env file")
		}
	}

	checker := helper.NewEnvVarChecker(
		constants.DatabaseUrl,
		constants.ServerPort,
		constants.Env,
		constants.Version,
		constants.WhatsappBusinessApiToken,
		constants.TelegramApiToken,
		constants.TelegramApiBotUrl,
		constants.GoogleClientId,
		constants.GoogleClientSecret,
		constants.JwtKey,
	)

	if envCheckerErr := checker.Check(); envCheckerErr != nil {
		log.Fatal(envCheckerErr.Error())
	}

	logger, loggerShutdown := logger.Create()

	db := database.New(logger)
	modules := initModules(db, logger)

	shudownFn := gracefulShutdown(loggerShutdown)

	return &App{
		Modules:  modules,
		Shutdown: shudownFn,
	}
}

func initModules(db *sql.DB, logger *slog.Logger) *Modules {
	return &Modules{
		Database: db,
		Logger:   logger,
		Telegram: nil,
		Whatsapp: nil,
	}
}

func gracefulShutdown(cleanups ...ShutdownFunc) ShutdownFunc {
	return func(ctx context.Context) error {
		var combinedErr error

		for i, cleanup := range cleanups {
			if cleanup == nil {
				continue
			}

			if err := cleanup(ctx); err != nil {
				slog.Error("Cleanup task failed", "task_index", i, "error", err)
				if combinedErr == nil {
					combinedErr = err
				} else {
					combinedErr = fmt.Errorf("%v; %w", combinedErr, err)
				}
			}
		}

		return combinedErr
	}
}
