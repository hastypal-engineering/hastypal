package controller

import (
	"log/slog"
	"net/http"
	"strconv"

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

		renderer := vendor.NewTemplRenderer(ctx, http.StatusOK, html.Error(&web.ErrorDTO{}))

		ctx.Render(http.StatusOK, renderer)

		return

		var req web.GetServicesReq
		if err := ctx.ShouldBindUri(&req); err != nil {
			c.logger.Error("Error binding GetServicesReq query params", "trace_id", traceID, "error", eris.ToString(err, true))

			renderer = vendor.NewTemplRenderer(ctx, http.StatusOK, html.Error(&web.ErrorDTO{}))

			ctx.Render(http.StatusOK, renderer)

			return
		}

		dto, err := c.service.ShowServices(ctx, req)
		if err != nil {
			c.logger.Error("Error showing services", "trace_id", traceID, "public_id", req.BusinessPublicID, "error", eris.ToString(err, true))

			renderer = vendor.NewTemplRenderer(ctx, http.StatusOK, html.Error(&web.ErrorDTO{}))

			ctx.Render(http.StatusOK, renderer)

			return
		}

		renderer = vendor.NewTemplRenderer(ctx, http.StatusOK, html.Step1(dto))

		ctx.Render(http.StatusOK, renderer)
	}
}

func (c *WebController) PostStep1() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.Value(middleware.TraceIDKey)

		publicID := ctx.Param("publicID")

		rawServiceID := ctx.PostForm("service")

		if rawServiceID == "" {
			c.logger.Error("Error storing selected service, serviceID missing", "trace_id", traceID, "public_id", publicID)

			renderer := vendor.NewTemplRenderer(ctx, http.StatusOK, html.Error(&web.ErrorDTO{}))

			ctx.Render(http.StatusOK, renderer)

			return
		}

		sessionID := ctx.PostForm("sessionID")

		if sessionID == "" {
			c.logger.Error("Error storing selected service, sessionID missing", "trace_id", traceID, "public_id", publicID)

			renderer := vendor.NewTemplRenderer(ctx, http.StatusOK, html.Error(&web.ErrorDTO{}))

			ctx.Render(http.StatusOK, renderer)

			return
		}

		serviceID, err := strconv.Atoi(rawServiceID)
		if err != nil {
			c.logger.Error("Error storing selected service while parsing rawServiceID", "trace_id", traceID, "public_id", publicID, "error", eris.ToString(err, true))

			renderer := vendor.NewTemplRenderer(ctx, http.StatusOK, html.Error(&web.ErrorDTO{}))

			ctx.Render(http.StatusOK, renderer)

			return
		}

		dto := &web.BookingPatchDTO{
			SessionID: sessionID,
			ServiceID: serviceID,
		}

		if err := c.service.StoreService(ctx, dto); err != nil {
			c.logger.Error("Error storing selected service", "trace_id", traceID, "public_id", publicID, "error", eris.ToString(err, true))

			renderer := vendor.NewTemplRenderer(ctx, http.StatusOK, html.Error(&web.ErrorDTO{}))

			ctx.Render(http.StatusOK, renderer)

			return
		}

		ctx.Redirect(203, "")
	}
}
