package booking

import (
	"context"
	"database/sql"
	"time"

	"github.com/rotisserie/eris"
)

type SessionRepository interface {
	Save(ctx context.Context, session *Session) error
	Update(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, sessionID string) (*Session, error)
	GetByDate(ctx context.Context, date time.Time) ([]*Session, error)
	GetByHour(ctx context.Context, date time.Time) (*Session, error)
}

type PgSessionRepository struct {
	connection *sql.DB
}

func NewPgSessionRepository(connection *sql.DB) *PgSessionRepository {
	return &PgSessionRepository{
		connection: connection,
	}
}

func (r *PgSessionRepository) Save(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO ha_booking_session (
			habs_id,
			habs_business_id,
			habs_service_id,
			habs_date,
			habs_hour,
			habs_ttl,
			habs_date_add,
			habs_date_upd
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err := r.connection.ExecContext(
		ctx,
		query,
		session.ID,
		session.BusinessID,
		session.ServiceID,
		session.Date,
		session.Hour,
		session.TTL,
		session.DateAdd,
		session.DateUpd,
	)

	if err != nil {
		return eris.Wrap(err, "Error saving session")
	}

	return nil
}

func (r *PgSessionRepository) GetByID(ctx context.Context, sessionID string) (*Session, error) {
	return nil, nil
}

func (r *PgSessionRepository) Update(ctx context.Context, session *Session) error {
	return nil
}

func (r *PgSessionRepository) GetByDate(ctx context.Context, date time.Time) ([]*Session, error) {
	return nil, nil
}

func (r *PgSessionRepository) GetByHour(ctx context.Context, date time.Time) (*Session, error) {
	return nil, nil
}
