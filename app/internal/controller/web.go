package controller

import (
	"log/slog"
	"net/http"

	"github.com/adriein/hastypal/internal/web"
	"github.com/adriein/hastypal/pkg/vendor"
	"github.com/adriein/hastypal/ui/html"
	"github.com/gin-gonic/gin"
)

type WebController struct {
	logger  *slog.Logger
	service web.WebService
}

func NewWebController(logger *slog.Logger, service web.WebService) *WebController {
	return &WebController{
		logger:  logger,
		service: service,
	}
}

func (c *WebController) GetStep1() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		renderer := vendor.NewTemplRenderer(ctx, http.StatusOK, html.Step1())

		ctx.Render(http.StatusOK, renderer)
	}
}

func (c *WebController) PostStep1() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{})
	}
}
