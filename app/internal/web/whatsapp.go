package web

import (
	"log/slog"
	"net/http"

	"github.com/adriein/hastypal/internal/whatsapp"
	"github.com/adriein/hastypal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
)

type WhatsappController struct {
	logger  *slog.Logger
	service whatsapp.WhatsappService
}

func NewWhatsappController(logger *slog.Logger, service whatsapp.WhatsappService) *WhatsappController {
	return &WhatsappController{
		logger:  logger,
		service: service,
	}
}

func (c *WhatsappController) Post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.Value(middleware.TraceIDKey)

		var update whatsapp.WhatsappUpdate

		if err := ctx.ShouldBindJSON(&update); err != nil {
			c.logger.Error("Error binding request to Whatsapp update struct", "trace_id", traceID, "error", eris.ToString(err, true))

			ctx.JSON(http.StatusInternalServerError, gin.H{})

			return
		}

		if err := c.service.HandleMessage(ctx, update); err != nil {
			c.logger.Error("Error handling Whatsapp update", "trace_id", traceID, "error", eris.ToString(err, true))

			ctx.JSON(http.StatusInternalServerError, gin.H{})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}
