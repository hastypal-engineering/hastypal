package booking

import (
	"context"
	"database/sql"
)

type BookingRepository interface {
	Save(ctx context.Context, booking *Booking) error
}

type PgBookingRepository struct {
	connection *sql.DB
}

func NewPgBookingRepository(connection *sql.DB) *PgBookingRepository {
	return &PgBookingRepository{
		connection: connection,
	}
}

func (r *PgBookingRepository) Save(ctx context.Context, booking *Booking) error {
	return nil
}
