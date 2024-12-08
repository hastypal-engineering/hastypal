package repository

import (
	"database/sql"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strconv"
	"strings"
)

type PgBookingSessionRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgBookingSessionRepository(connection *sql.DB) *PgBookingSessionRepository {
	transformer, _ := helper.NewCriteriaToSqlService(&types.BookingSession{})

	return &PgBookingSessionRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgBookingSessionRepository) Find(criteria types.Criteria) ([]types.BookingSession, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, exception.
			New(err.Error()).
			Trace("r.transformer.Transform", "pg-booking-session-repository.go")
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, exception.
			New(queryErr.Error()).
			Trace("r.connection.Query", "pg-booking-session-repository.go")
	}

	defer rows.Close()

	var (
		id          string
		business_id string
		chat_id     int
		service_id  string
		date        string
		hour        string
		created_at  string
		updated_at  string
		ttl         int64
	)

	var results []types.BookingSession

	for rows.Next() {
		if scanErr := rows.Scan(
			&id,
			&business_id,
			&chat_id,
			&service_id,
			&date,
			&hour,
			&created_at,
			&updated_at,
			&ttl,
		); scanErr != nil {
			return nil, exception.
				New(scanErr.Error()).
				Trace("rows.Scan", "pg-booking-session-repository.go")
		}

		results = append(results, types.BookingSession{
			Id:         id,
			BusinessId: business_id,
			ChatId:     chat_id,
			ServiceId:  service_id,
			Date:       date,
			Hour:       hour,
			CreatedAt:  created_at,
			UpdatedAt:  updated_at,
			Ttl:        ttl,
		})
	}

	return results, nil
}

func (r *PgBookingSessionRepository) FindOne(criteria types.Criteria) (types.BookingSession, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types.BookingSession{}, exception.
			New(err.Error()).
			Trace("r.transformer.Transform", "pg-booking-session-repository.go")
	}

	var (
		id          string
		business_id string
		chat_id     int
		service_id  string
		date        string
		hour        string
		created_at  string
		updated_at  string
		ttl         int64
	)

	if scanErr := r.connection.QueryRow(query).Scan(
		&id,
		&business_id,
		&chat_id,
		&service_id,
		&date,
		&hour,
		&created_at,
		&updated_at,
		&ttl,
	); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return types.BookingSession{}, exception.
				New(scanErr.Error()).
				Trace("r.connection.QueryRow", "pg-booking-session-repository.go").
				WithValues([]string{query}).
				Domain()
		}

		return types.BookingSession{}, exception.
			New(scanErr.Error()).
			Trace("r.connection.QueryRow", "pg-booking-session-repository.go").
			WithValues([]string{query})
	}

	return types.BookingSession{
		Id:         id,
		BusinessId: business_id,
		ChatId:     chat_id,
		ServiceId:  service_id,
		Date:       date,
		Hour:       hour,
		CreatedAt:  created_at,
		UpdatedAt:  updated_at,
		Ttl:        ttl,
	}, nil
}

func (r *PgBookingSessionRepository) Save(entity types.BookingSession) error {
	var query strings.Builder

	query.WriteString(`INSERT INTO booking_session `)
	query.WriteString(`(id, business_id, chat_id, service_id, date, hour, created_at, updated_at, ttl) `)
	query.WriteString(`VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`)

	_, err := r.connection.Exec(
		query.String(),
		entity.Id,
		entity.BusinessId,
		entity.ChatId,
		entity.ServiceId,
		entity.Date,
		entity.Hour,
		entity.CreatedAt,
		entity.UpdatedAt,
		entity.Ttl,
	)

	if err != nil {
		return exception.
			New(err.Error()).
			Trace("r.connection.Exec", "pg-booking-session-repository.go").
			WithValues([]string{
				query.String(),
				entity.Id,
				entity.BusinessId,
				strconv.Itoa(entity.ChatId),
			})
	}

	return nil
}

func (r *PgBookingSessionRepository) Update(entity types.BookingSession) error {
	var query strings.Builder

	query.WriteString(`UPDATE booking_session `)
	query.WriteString(`SET business_id = $2, chat_id = $3, service_id = $4, date = $5, hour = $6, `)
	query.WriteString(`created_at = $7, updated_at = $8, ttl = $9 WHERE id = $1;`)

	_, err := r.connection.Exec(
		query.String(),
		entity.Id,
		entity.BusinessId,
		entity.ChatId,
		entity.ServiceId,
		entity.Date,
		entity.Hour,
		entity.CreatedAt,
		entity.UpdatedAt,
		entity.Ttl,
	)

	if err != nil {
		return exception.
			New(err.Error()).
			Trace("r.connection.Exec", "pg-booking-session-repository.go").
			WithValues([]string{
				query.String(),
				entity.Id,
				entity.BusinessId,
				strconv.Itoa(entity.ChatId),
			})
	}

	return nil
}
