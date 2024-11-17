package main

import (
	"database/sql"
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/business"
	"github.com/adriein/hastypal/internal/hastypal/google"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/repository"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	"log"
	"net/http"
	"os"

	"github.com/adriein/hastypal/internal/hastypal/handler"
	"github.com/adriein/hastypal/internal/hastypal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
		constants.GoogleClientId,
		constants.GoogleClientSecret,
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

	api.Route("POST /telegram-webhook", constructTelegramWebhookHandler(api, database))

	api.Route("GET /business/google-auth", constructGoogleAuthHandler(api))
	api.Route("GET /business/google-auth-callback", constructGoogleAuthCallbackHandler(api, database))

	api.Route("POST /business", constructCreateBusinessHandler(api, database))

	api.Start()

	defer database.Close()
}

func constructTelegramWebhookHandler(api *server.HastypalApiServer, database *sql.DB) http.HandlerFunc {
	googleApi := service.NewGoogleApi()
	businessRepository := repository.NewPgBusinessRepository(database)
	sessionRepository := repository.NewPgBookingSessionRepository(database)
	notificationRepository := repository.NewPgTelegramNotificationRepository(database)
	bookingRepository := repository.NewPgBookingRepository(database)
	googleTokenRepository := repository.NewPgGoogleTokenRepository(database)

	bot := service.NewTelegramBot(os.Getenv(constants.TelegramApiBotUrl), os.Getenv(constants.TelegramApiToken))

	startCommandHandler := service.NewTelegramStartCommandService(bot, sessionRepository)
	datesCommandHandler := service.NewTelegramDatesCommandService(bot, sessionRepository)
	hoursCommandHandler := service.NewTelegramHoursCommandService(bot, sessionRepository)
	confirmationCommandHandler := service.NewTelegramConfirmationCommandService(bot, sessionRepository)
	finishCommandHandler := service.NewTelegramFinishCommandService(
		bot,
		googleApi,
		sessionRepository,
		notificationRepository,
		bookingRepository,
		googleTokenRepository,
	)

	webhookService := service.NewTelegramWebhookService(
		businessRepository,
		startCommandHandler,
		datesCommandHandler,
		hoursCommandHandler,
		confirmationCommandHandler,
		finishCommandHandler,
	)

	controller := handler.NewTelegramWebhookHandler(webhookService)

	return api.NewHandler(controller.Handler)
}

func constructGoogleAuthHandler(api *server.HastypalApiServer) http.HandlerFunc {
	googleApi := service.NewGoogleApi()

	authGoogleService := google.NewAuthGoogleService(googleApi)

	controller := google.NewGoogleAuthHandler(authGoogleService)

	return api.NewHandler(controller.Handler)
}

func constructGoogleAuthCallbackHandler(api *server.HastypalApiServer, database *sql.DB) http.HandlerFunc {
	googleApi := service.NewGoogleApi()
	googleTokenRepository := repository.NewPgGoogleTokenRepository(database)

	callbackService := google.NewAuthCallbackGoogleService(googleTokenRepository, googleApi)

	controller := google.NewGoogleAuthCallbackHandler(callbackService)

	return api.NewHandler(controller.Handler)
}

func constructCreateBusinessHandler(api *server.HastypalApiServer, database *sql.DB) http.HandlerFunc {
	businessRepository := repository.NewPgBusinessRepository(database)

	createBusinessService := business.NewCreateBusinessService(businessRepository)

	controller := business.NewCreateBusinessHandler(createBusinessService)

	return api.NewHandler(controller.Handler)
}
