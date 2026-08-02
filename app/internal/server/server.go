package server

import (
	"fmt"
	"os"

	"github.com/adriein/hastypal/internal"
	"github.com/adriein/hastypal/internal/web"
	"github.com/adriein/hastypal/pkg/constants"
	"github.com/adriein/hastypal/pkg/middleware"
	"github.com/adriein/hastypal/pkg/vendor"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rotisserie/eris"
)

type Server struct {
	gin       *gin.Engine
	validator *validator.Validate
}

func New(app *internal.App) *Server {
	logger := app.Modules.Logger
	engine := gin.New()

	ginHtmlRenderer := engine.HTMLRender

	engine.HTMLRender = &vendor.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	if os.Getenv(constants.Env) == constants.Pro {
		gin.SetMode(gin.ReleaseMode)
	}

	engine.SetTrustedProxies(nil)

	engine.Use(middleware.Error(), gin.Logger(), gin.Recovery(), middleware.Tracer(), middleware.TimeZone())

	server := &Server{
		gin:       engine,
		validator: validator.New(),
	}

	server.routeSetup()

	port := os.Getenv(constants.ServerPort)

	if ginErr := engine.Run(port); ginErr != nil {
		err := eris.Wrap(ginErr, "Error starting HTTP server")

		logger.Error(eris.ToString(err, true))
		os.Exit(1)
	}

	logger.Info("Starting the TibiaChar at " + port)

	return server
}

func (s *Server) routeSetup() {
	//HEALTH CHECK
	s.gin.GET("/health", web.NewHealthController().Get())

	cwd, _ := os.Getwd()

	//STATIC
	s.gin.Static("/ui/static", fmt.Sprintf("%s/ui/static", cwd))

	//TODO: setup the routes again

	/*
		api.Route("POST /telegram-webhook", constructTelegramWebhookHandler(api, database))

		api.Route("GET /business/google-auth", constructGoogleAuthHandler(api))
		api.Route("GET /business/google-auth-callback", constructGoogleAuthCallbackHandler(api, database))
		api.Route("POST /business", constructCreateBusinessHandler(api, database))
		api.Route("POST /business/login", constructLoginBusinessHandler(api, database))

		api.Route("GET /notification/send", constructSendNotificationHandler(api, database))
	*/
}
