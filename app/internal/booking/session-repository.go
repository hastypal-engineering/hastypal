package booking

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/adriein/hastypal/database"
	"github.com/adriein/hastypal/pkg/helper/conversion"
	"github.com/rotisserie/eris"
)

type SessionRepository interface {
	Save(ctx context.Context, session *Session) error
	Update(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, sessionID string) (*Session, error)
	GetByDate(ctx context.Context, businessID int, date time.Time) ([]*Session, error)
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
			habs_start_time,
			habs_end_time,
			habs_ttl,
			habs_date_add,
			habs_date_upd
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`

	_, err := r.connection.ExecContext(
		ctx,
		query,
		session.ID,
		session.BusinessID,
		session.ServiceID,
		conversion.DateToDB(session.Date),
		conversion.HourToDB(session.Interval.Start),
		conversion.HourToDB(session.Interval.End),
		session.TTL.Milliseconds(),
		session.DateAdd,
		session.DateUpd,
	)
	if err != nil {
		return eris.Wrap(err, "Error saving session")
	}

	return nil
}

func (r *PgSessionRepository) GetByID(ctx context.Context, sessionID string) (*Session, error) {
	query := `
		SELECT
			habs_id,
			habs_business_id,
			habs_service_id,
			habs_date,
			habs_start_time,
			habs_end_time,
			habs_ttl,
			habs_date_add,
			habs_date_upd
		FROM
			ha_booking_session
		WHERE
			habs_id = $1;
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	session, err := scanSession(r.connection.QueryRowContext(ctxTimeout, query, sessionID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, eris.Wrap(err, "Session not found by ID")
		}

		return nil, eris.Wrap(err, "Failed to query session by ID")
	}

	return session, nil
}

func (r *PgSessionRepository) Update(ctx context.Context, session *Session) error {
	query := `
		UPDATE ha_booking_session
		SET
			habs_service_id = $2,
			habs_date = $3,
			habs_start_time = $4,
			habs_end_time = $5,
			habs_ttl = $6,
			habs_date_upd = $7
		WHERE
			habs_id = $1;
	`

	_, err := r.connection.ExecContext(
		ctx,
		query,
		session.ID,
		session.ServiceID,
		conversion.DateToDB(session.Date),
		conversion.HourToDB(session.Interval.Start),
		conversion.HourToDB(session.Interval.End),
		session.TTL.Milliseconds(),
		session.DateUpd,
	)
	if err != nil {
		return eris.Wrap(err, "Error updating session")
	}

	return nil
}

func (r *PgSessionRepository) GetByDate(ctx context.Context, businessID int, date time.Time) ([]*Session, error) {
	query := `
		SELECT
			habs_id,
			habs_business_id,
			habs_service_id,
			habs_date,
			habs_start_time,
			habs_end_time,
			habs_ttl,
			habs_date_add,
			habs_date_upd
		FROM
			ha_booking_session
		WHERE
			habs_date >= $1
		AND
			habs_business_id = $2
		ORDER BY
			habs_start_time ASC;
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rows, err := r.connection.QueryContext(ctxTimeout, query, date, businessID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to query sessions by date")
	}

	database.CloseRowsSafely(rows, &err)

	sessions := []*Session{}

	for rows.Next() {
		session, err := scanSession(rows)
		if err != nil {
			return nil, eris.Wrap(err, "Failed to scan a session row while querying by date")
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, eris.Wrap(err, "Failed while iterating the sessions queried by date")
	}

	return sessions, nil
}

func (r *PgSessionRepository) GetByHour(ctx context.Context, date time.Time) (*Session, error) {
	return nil, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanSession(row scanner) (*Session, error) {
	session := &Session{}

	var (
		serviceID sql.NullInt64
		date      sql.NullTime
		start     sql.NullTime
		end       sql.NullTime
		ttl       sql.NullInt64
	)

	if err := row.Scan(
		&session.ID,
		&session.BusinessID,
		&serviceID,
		&date,
		&start,
		&end,
		&ttl,
		&session.DateAdd,
		&session.DateUpd,
	); err != nil {
		return nil, err
	}

	session.ServiceID = int(serviceID.Int64)
	session.TTL = time.Duration(ttl.Int64) * time.Millisecond

	if date.Valid {
		session.Date = date.Time
	}

	if start.Valid && end.Valid {
		session.Interval.Start = conversion.HourFromDB(start.Time)
		session.Interval.End = conversion.HourFromDB(end.Time)

		return session, nil
	}

	session.Interval = &SessionInterval{}

	return session, nil
}
