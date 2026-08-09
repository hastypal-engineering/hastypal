package booking

import (
	"context"
	"database/sql"

	"github.com/rotisserie/eris"
)

type SessionRepository interface {
	Save(ctx context.Context, session *Session) error
	Update(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, sessionID string) (*Session, error)
}

type PgSessionRepository struct {
	connection *sql.DB
}

func NewPgSessionRepository(connection *sql.DB) *PgSessionRepository {
	return &PgSessionRepository{
		connection: connection,
	}
}

func (r *PgSessionRepository) Save(ctx context.Context, session Session) error {
	query := `
		INSERT INTO tc_currency_rates (
			tcr_usd,
			tcr_eur,
			tcr_aud,
			tcr_gbp,
			tcr_pln,
			tcr_brl,
			tcr_date_upd
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	_, err := r.connection.ExecContext(
		ctx,
		query,
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
