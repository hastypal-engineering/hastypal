package main

import (
	"database/sql"
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/constants"
	"github.com/adriein/hastypal/internal/hastypal/handler"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/repository"
	"github.com/adriein/hastypal/internal/hastypal/server"
	"github.com/adriein/hastypal/internal/hastypal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	dotenvErr := godotenv.Load()

	if dotenvErr != nil && os.Getenv(constants.Env) != constants.Production {
		log.Fatal("Error loading .env file")
	}

	checker := helper.NewEnvVarChecker(
		constants.DatabaseUser,
		constants.DatabasePassword,
		constants.DatabaseName,
		constants.ServerPort,
		constants.WhatsappBusinessApiToken,
		constants.TelegramApiToken,
		constants.TelegramApiBotUrl,
	)

	if envCheckerErr := checker.Check(); envCheckerErr != nil {
		log.Fatal(envCheckerErr.Error())
	}

	api, newServerErr := server.New(os.Getenv(constants.ServerPort))

	if newServerErr != nil {
		log.Fatal(newServerErr.Error())
	}

	databaseDsn := fmt.Sprintf(
		"postgresql://%s:%s@localhost:5432/%s?sslmode=disable",
		os.Getenv(constants.DatabaseUser),
		os.Getenv(constants.DatabasePassword),
		os.Getenv(constants.DatabaseName),
	)

	database, dbConnErr := sql.Open("postgres", databaseDsn)

	if dbConnErr != nil {
		log.Fatal(dbConnErr.Error())
	}

	api.Route("POST /bot-setup", constructBotSetupHandler(api, database))
	api.Route("POST /telegram-webhook", constructTelegramWebhookHandler(api, database))

	api.Start()

	defer database.Close()
}

func constructBotSetupHandler(api *server.HastypalApiServer, database *sql.DB) http.HandlerFunc {
	bot := service.NewTelegramBot(os.Getenv(constants.TelegramApiBotUrl), os.Getenv(constants.TelegramApiToken))

	setupBotService := service.NewSetupTelegramBotService(bot)

	controller := handler.NewSetupTelegramBotHandler(setupBotService)

	return api.NewHandler(controller.Handler)
}

func constructTelegramWebhookHandler(api *server.HastypalApiServer, database *sql.DB) http.HandlerFunc {
	bot := service.NewTelegramBot(os.Getenv(constants.TelegramApiBotUrl), os.Getenv(constants.TelegramApiToken))
	businessRepository := repository.NewPgBusinessRepository(database)

	startCommandHandler := service.NewTelegramStartCommandService(bot)
	datesCommandHandler := service.NewTelegramDatesCommandService(bot)
	hoursCommandHandler := service.NewTelegramHoursCommandService(bot)
	confirmationCommandHandler := service.NewTelegramConfirmationCommandService(bot)

	webhookService := service.NewTelegramWebhookService(
		businessRepository,
		startCommandHandler,
		datesCommandHandler,
		hoursCommandHandler,
		confirmationCommandHandler,
	)

	controller := handler.NewTelegramWebhookHandler(webhookService)

	return api.NewHandler(controller.Handler)
}
