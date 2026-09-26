package business

import "log/slog"

type HolidayAPI struct {
	logger *slog.Logger
}

func NewHolidayAPI(logger *slog.Logger) *HolidayAPI {
	return &HolidayAPI{
		logger: logger,
	}
}
