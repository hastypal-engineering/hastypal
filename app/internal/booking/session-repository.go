package booking

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/adriein/hastypal/database"
	"github.com/adriein/hastypal/pkg/helper"
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
			habs_slot_index,
			habs_client_name,
			habs_client_email,
			habs_client_phone,
			habs_notes,
			habs_date_add,
			habs_date_upd
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
	`

	_, err := r.connection.ExecContext(
		ctx,
		query,
		session.ID,
		session.BusinessID,
		session.ServiceID,
		helper.DateToDB(session.Date),
		helper.HourToDB(session.Interval.Start),
		helper.HourToDB(session.Interval.End),
		session.TTL.Milliseconds(),
		session.SlotIndex,
		nullableString(session.ClientName),
		nullableString(session.ClientEmail),
		nullableString(session.ClientPhone),
		nullableString(session.Notes),
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
			habs_slot_index,
			habs_client_name,
			habs_client_email,
			habs_client_phone,
			habs_notes,
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
			habs_slot_index = $7,
			habs_client_name = $8,
			habs_client_email = $9,
			habs_client_phone = $10,
			habs_notes = $11,
			habs_date_upd = $12
		WHERE
			habs_id = $1;
	`

	_, err := r.connection.ExecContext(
		ctx,
		query,
		session.ID,
		session.ServiceID,
		helper.DateToDB(session.Date),
		helper.HourToDB(session.Interval.Start),
		helper.HourToDB(session.Interval.End),
		session.TTL.Milliseconds(),
		session.SlotIndex,
		nullableString(session.ClientName),
		nullableString(session.ClientEmail),
		nullableString(session.ClientPhone),
		nullableString(session.Notes),
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
			habs_slot_index,
			habs_client_name,
			habs_client_email,
			habs_client_phone,
			habs_notes,
			habs_date_add,
			habs_date_upd
		FROM
			ha_booking_session
		WHERE
			habs_date = $1
		AND
			habs_business_id = $2
		ORDER BY
			habs_start_time ASC;
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rows, err := r.connection.QueryContext(ctxTimeout, query, date.Format(time.DateOnly), businessID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to query sessions by date")
	}

	defer database.CloseRowsSafely(rows, &err)

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

func nullableString(s string) any {
	if s == "" {
		return nil
	}

	return s
}

func scanSession(row scanner) (*Session, error) {
	session := &Session{
		Interval: &SessionInterval{},
	}

	var (
		serviceID   sql.NullInt64
		date        sql.NullTime
		start       sql.NullTime
		end         sql.NullTime
		ttl         sql.NullInt64
		slotIndex   sql.NullInt64
		clientName  sql.NullString
		clientEmail sql.NullString
		clientPhone sql.NullString
		notes       sql.NullString
	)

	if err := row.Scan(
		&session.ID,
		&session.BusinessID,
		&serviceID,
		&date,
		&start,
		&end,
		&ttl,
		&slotIndex,
		&clientName,
		&clientEmail,
		&clientPhone,
		&notes,
		&session.DateAdd,
		&session.DateUpd,
	); err != nil {
		return nil, err
	}

	session.ServiceID = int(serviceID.Int64)
	session.TTL = time.Duration(ttl.Int64) * time.Millisecond

	if slotIndex.Valid {
		session.SlotIndex = int(slotIndex.Int64)
	}

	if clientName.Valid {
		session.ClientName = clientName.String
	}

	if clientEmail.Valid {
		session.ClientEmail = clientEmail.String
	}

	if clientPhone.Valid {
		session.ClientPhone = clientPhone.String
	}

	if notes.Valid {
		session.Notes = notes.String
	}

	if date.Valid {
		session.Date = date.Time
	}

	if start.Valid && end.Valid {
		session.Interval.Start = helper.HourFromDB(start.Time)
		session.Interval.End = helper.HourFromDB(end.Time)

		return session, nil
	}

	return session, nil
}
