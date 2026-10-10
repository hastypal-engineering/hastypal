package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/adriein/hastypal/internal/booking"
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

func redirectOnError(ctx *gin.Context, location string) {
	if ctx.GetHeader("HX-Request") == "true" {
		ctx.Header("HX-Redirect", location)
		ctx.Status(http.StatusOK)

		return
	}

	ctx.Redirect(http.StatusFound, location)
}

func sessionExpiredLocation(publicID string) string {
	return fmt.Sprintf("/booking/%s/session-expired", publicID)
}

func (c *WebController) GetStep1() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.Value(middleware.TraceIDKey)

		var req web.GetServicesReq
		if err := ctx.ShouldBindUri(&req); err != nil {
			c.logger.Error("Error binding GetServicesReq query params", "trace_id", traceID, "error", eris.ToString(err, true))

			redirectOnError(ctx, fmt.Sprintf("/booking/%s/error", req.BusinessPublicID))

			return
		}

		dto, err := c.service.ShowServices(ctx, req)
		if err != nil {
			if errors.Is(err, booking.BookingSessionExpired) {
				redirectOnError(ctx, sessionExpiredLocation(req.BusinessPublicID))

				return
			}

			c.logger.Error("Error showing services", "trace_id", traceID, "public_id", req.BusinessPublicID, "error", eris.ToString(err, true))

			redirectOnError(ctx, fmt.Sprintf("/booking/%s/error", req.BusinessPublicID))

			return
		}

		renderer := vendor.NewTemplRenderer(ctx, http.StatusOK, html.Step1(dto))

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

			redirectOnError(ctx, fmt.Sprintf("/booking/%s/error", publicID))

			return
		}

		sessionID := ctx.PostForm("sessionID")

		if sessionID == "" {
			c.logger.Error("Error storing selected service, sessionID missing", "trace_id", traceID, "public_id", publicID)

			redirectOnError(ctx, fmt.Sprintf("/booking/%s/error", publicID))

			return
		}

		serviceID, err := strconv.Atoi(rawServiceID)
		if err != nil {
			c.logger.Error("Error storing selected service while parsing rawServiceID", "trace_id", traceID, "public_id", publicID, "error", eris.ToString(err, true))

			redirectOnError(ctx, fmt.Sprintf("/booking/%s/error", publicID))

			return
		}

		dto := &web.BookingPatchDTO{
			SessionID: sessionID,
			ServiceID: serviceID,
		}

		if err := c.service.StoreService(ctx, dto); err != nil {
			if errors.Is(err, booking.BookingSessionExpired) {
				redirectOnError(ctx, sessionExpiredLocation(publicID))

				return
			}

			c.logger.Error("Error storing selected service", "trace_id", traceID, "public_id", publicID, "error", eris.ToString(err, true))

			redirectOnError(ctx, fmt.Sprintf("/booking/%s/error", publicID))

			return
		}

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/booking/%s/session/%s/step-2", publicID, sessionID))
	}
}

func (c *WebController) GetStep2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.Value(middleware.TraceIDKey)

		var req web.GetDatesReq
		if err := ctx.ShouldBindUri(&req); err != nil {
			c.logger.Error("Error showing dates while binding GetDatesReq uri params", "trace_id", traceID, "error", eris.ToString(err, true))

			redirectOnError(ctx, fmt.Sprintf("/booking/%s/error", req.BusinessPublicID))

			return
		}

		if err := ctx.ShouldBindQuery(&req); err != nil {
			c.logger.Error("Error showing dates while binding GetDatesReq query params", "trace_id", traceID, "error", eris.ToString(err, true))

			redirectOnError(ctx, fmt.Sprintf("/booking/%s/error", req.BusinessPublicID))

			return
		}

		dto, err := c.service.ShowDates(ctx, req)
		if err != nil {
			if errors.Is(err, booking.BookingSessionExpired) {
				redirectOnError(ctx, sessionExpiredLocation(req.BusinessPublicID))

				return
			}

			c.logger.Error("Error showing dates", "trace_id", traceID, "public_id", req.BusinessPublicID, "session_id", req.SessionID, "error", eris.ToString(err, true))

			redirectOnError(ctx, fmt.Sprintf("/booking/%s/error", req.BusinessPublicID))

			return
		}

		renderer := vendor.NewTemplRenderer(ctx, http.StatusOK, html.Step2(dto))

		ctx.Render(http.StatusOK, renderer)
	}
}

func (c *WebController) GetError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		renderer := vendor.NewTemplRenderer(ctx, http.StatusInternalServerError, html.Error(&web.ErrorDTO{}))

		ctx.Render(http.StatusInternalServerError, renderer)
	}
}

func (c *WebController) GetSessionExpired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		publicID := ctx.Param("publicID")

		renderer := vendor.NewTemplRenderer(
			ctx,
			http.StatusUnauthorized,
			html.SessionExpired(&web.ErrorDTO{BusinessPublicID: publicID}),
		)

		ctx.Render(http.StatusUnauthorized, renderer)
	}
}
