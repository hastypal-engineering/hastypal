package controller

import (
	"log/slog"
	"net/http"

	"github.com/adriein/hastypal/internal/web"
	"github.com/adriein/hastypal/pkg/middleware"
	"github.com/adriein/hastypal/pkg/vendor"
	"github.com/adriein/hastypal/ui/html"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
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
		traceID := ctx.Value(middleware.TraceIDKey)

		var req web.GetServicesReq

		if err := ctx.ShouldBindJSON(&req); err != nil {
			c.logger.Error("Error binding request to telegram update struct", "trace_id", traceID, "error", eris.ToString(err, true))

			return
		}

		dtos, err := c.service.ShowServices(ctx, req)

		if err != nil {
			c.logger.Error("Error binding request to telegram update struct", "trace_id", traceID, "error", eris.ToString(err, true))
		}

		renderer := vendor.NewTemplRenderer(ctx, http.StatusOK, html.Step1(dtos))

		ctx.Render(http.StatusOK, renderer)
	}
}

func (c *WebController) PostStep1() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}
