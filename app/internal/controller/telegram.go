package controller

import (
	"log/slog"
	"net/http"

	"github.com/adriein/hastypal/internal/telegram"
	"github.com/adriein/hastypal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
)

type TelegramController struct {
	logger  *slog.Logger
	service telegram.TelegramService
}

func NewTelegramController(logger *slog.Logger, service telegram.TelegramService) *TelegramController {
	return &TelegramController{
		logger:  logger,
		service: service,
	}
}

func (c *TelegramController) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.Value(middleware.TraceIDKey)

		var update telegram.TelegramUpdate

		if err := ctx.ShouldBindJSON(&update); err != nil {
			c.logger.Error("Error binding request to telegram update struct", "trace_id", traceID, "error", eris.ToString(err, true))

			ctx.JSON(http.StatusInternalServerError, gin.H{})

			return
		}

		if err := c.service.HandleMessage(ctx, update); err != nil {
			c.logger.Error("Error handling telegram update", "trace_id", traceID, "error", eris.ToString(err, true))

			ctx.JSON(http.StatusInternalServerError, gin.H{})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}
