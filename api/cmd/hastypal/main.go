package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adriein/hastypal/internal/hastypal/notification"
	"github.com/adriein/hastypal/internal/hastypal/shared/translation"

	"github.com/adriein/hastypal/internal/hastypal/business"
	"github.com/adriein/hastypal/internal/hastypal/google"
	"github.com/adriein/hastypal/internal/hastypal/server"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/repository"
	"github.com/adriein/hastypal/internal/hastypal/shared/service"
	"github.com/adriein/hastypal/internal/hastypal/telegram"
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
		constants.JwtKey,
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

	// To define middlewares:
	// cronMiddlewares := middleware.NewMiddlewareChain(
	// 	middleware.NewAuthMiddleWare,
	// )

	api.Route("POST /telegram-webhook", constructTelegramWebhookHandler(api, database))

	api.Route("GET /business/google-auth", constructGoogleAuthHandler(api))
	api.Route("GET /business/google-auth-callback", constructGoogleAuthCallbackHandler(api, database))
	api.Route("POST /business", constructCreateBusinessHandler(api, database))
	api.Route("POST /business/login", constructLoginBusinessHandler(api, database))

	// To apply auth middleware in an endpoint:
	// api.Route("VERB /endpoint", cronMiddlewares.ApplyOn(handlerConstructor))

	api.Route("GET /notification/send", constructSendNotificationHandler(api, database))

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
	translations := translation.New()

	startCommandService := telegram.NewStartCommandTelegramService(bot, sessionRepository, businessRepository)
	datesCommandService := telegram.NewPickDateCommandTelegramService(bot, sessionRepository, translations)
	hoursCommandService := telegram.NewPickHourCommandTelegramService(bot, sessionRepository, translations)
	confirmationCommandService := telegram.NewConfirmationCommandTelegramService(bot, sessionRepository, translations)
	finishCommandService := telegram.NewFinishCommandTelegramService(
		bot,
		googleApi,
		sessionRepository,
		notificationRepository,
		bookingRepository,
		googleTokenRepository,
		businessRepository,
	)

	webhookService := telegram.NewNotificationWebhookTelegramService(
		startCommandService,
		datesCommandService,
		hoursCommandService,
		confirmationCommandService,
		finishCommandService,
	)

	controller := telegram.NewNotificationWebhookTelegramHandler(webhookService)

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

func constructLoginBusinessHandler(api *server.HastypalApiServer, database *sql.DB) http.HandlerFunc {
	businessRepository := repository.NewPgBusinessRepository(database)

	loginBusinessService := business.NewLoginBusinessService(businessRepository)

	controller := business.NewLoginBusinessHandler(loginBusinessService)

	return api.NewHandler(controller.Handler)
}

func constructSendNotificationHandler(api *server.HastypalApiServer, database *sql.DB) http.HandlerFunc {
	bot := service.NewTelegramBot(os.Getenv(constants.TelegramApiBotUrl), os.Getenv(constants.TelegramApiToken))
	notificationRepository := repository.NewPgTelegramNotificationRepository(database)
	translations := translation.New()

	sendNotificationService := notification.NewSendNotificationService(notificationRepository, bot, translations)

	controller := notification.NewSendNotificationHandler(sendNotificationService)

	return api.NewHandler(controller.Handler)
}
