// Package seed
package seed

import (
	"database/sql"
	"log/slog"
)

type SeedService interface {
	Run(connection *sql.DB, logger *slog.Logger) error
}

type Service struct{}

func Run(connection *sql.DB, logger *slog.Logger) error {
	return nil
}
