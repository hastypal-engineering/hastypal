package server

import (
	"fmt"
	"os"

	"github.com/adriein/hastypal/internal"
	"github.com/adriein/hastypal/internal/controller"
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

	engine.Use(gin.Logger(), gin.Recovery(), middleware.Tracer(), middleware.TimeZone())

	server := &Server{
		gin:       engine,
		validator: validator.New(),
	}

	server.routeSetup(app)

	port := os.Getenv(constants.ServerPort)

	if ginErr := engine.Run(port); ginErr != nil {
		err := eris.Wrap(ginErr, "Error starting HTTP server")

		logger.Error(eris.ToString(err, true))
		os.Exit(1)
	}

	logger.Info("Starting the TibiaChar at " + port)

	return server
}

func (s *Server) routeSetup(app *internal.App) {
	//HEALTH CHECK
	s.gin.GET("/health", controller.NewHealthController().Get())

	//TELEGRAM WEBHOOK

	s.gin.POST("/telegram-webhook", s.webhookController(app).Post())

	s.gin.GET("/booking/step-1", s.webController(app).GetStep1())

	cwd, _ := os.Getwd()

	//STATIC
	s.gin.Static("/ui/static", fmt.Sprintf("%s/ui/static", cwd))

	/*
		api.Route("GET /business/google-auth", constructGoogleAuthHandler(api))
		api.Route("GET /business/google-auth-callback", constructGoogleAuthCallbackHandler(api, database))
		api.Route("POST /business", constructCreateBusinessHandler(api, database))
		api.Route("POST /business/login", constructLoginBusinessHandler(api, database))

		api.Route("GET /notification/send", constructSendNotificationHandler(api, database))
	*/
}

func (s *Server) webhookController(app *internal.App) *controller.TelegramController {
	logger := app.Modules.Logger
	service := app.Modules.Telegram

	return controller.NewTelegramController(logger, service)
}

func (s *Server) webController(app *internal.App) *controller.WebController {
	logger := app.Modules.Logger
	service := app.Modules.Web

	return controller.NewWebController(logger, service)
}
