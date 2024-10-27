package repository

import (
	"database/sql"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
	"strconv"
	"strings"
)

type PgBookingSessionRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgBookingSessionRepository(connection *sql.DB) *PgBookingSessionRepository {
	transformer := helper.NewCriteriaToSqlService("booking_session")

	return &PgBookingSessionRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgBookingSessionRepository) Find(criteria types.Criteria) ([]types.BookingSession, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, types.ApiError{
			Msg:      err.Error(),
			Function: "Find -> r.transformer.Transform()",
			File:     "pg-booking-session-repository.go",
		}
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, types.ApiError{
			Msg:      queryErr.Error(),
			Function: "Find -> r.connection.Query()",
			File:     "pg-booking-session-repository.go",
			Values:   []string{query},
		}
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
			&ttl,
		); scanErr != nil {
			return nil, types.ApiError{
				Msg:      scanErr.Error(),
				Function: "Find -> rows.Scan()",
				File:     "pg-booking-session-repository.go",
			}
		}

		results = append(results, types.BookingSession{
			Id:         id,
			BusinessId: business_id,
			ChatId:     chat_id,
			ServiceId:  service_id,
			Date:       date,
			Hour:       hour,
			CreatedAt:  created_at,
			Ttl:        ttl,
		})
	}

	return results, nil
}

func (r *PgBookingSessionRepository) FindOne(criteria types.Criteria) (types.BookingSession, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types.BookingSession{}, types.ApiError{
			Msg:      err.Error(),
			Function: "FindOne -> r.transformer.Transform()",
			File:     "pg-booking-session-repository.go",
		}
	}

	var (
		id          string
		business_id string
		chat_id     int
		service_id  string
		date        string
		hour        string
		created_at  string
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
		&ttl,
	); scanErr != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return types.BookingSession{}, types.ApiError{
				Msg:      "Entity Business not found",
				Function: "FindOne -> rows.Scan()",
				File:     "pg-booking-session-repository.go",
				Values:   []string{query},
				Domain:   true,
			}
		}

		return types.BookingSession{}, types.ApiError{
			Msg:      scanErr.Error(),
			Function: "FindOne -> rows.Scan()",
			File:     "pg-booking-session-repository.go",
			Values:   []string{query},
		}
	}

	return types.BookingSession{
		Id:         id,
		BusinessId: business_id,
		ChatId:     chat_id,
		ServiceId:  service_id,
		Date:       date,
		Hour:       hour,
		CreatedAt:  created_at,
		Ttl:        ttl,
	}, nil
}

func (r *PgBookingSessionRepository) Save(entity types.BookingSession) error {
	var query strings.Builder

	query.WriteString(`INSERT INTO booking_session `)
	query.WriteString(`(id, business_id, chat_id, service_id, date, hour, created_at, ttl) `)
	query.WriteString(`VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`)

	_, err := r.connection.Exec(
		query.String(),
		entity.Id,
		entity.BusinessId,
		entity.ChatId,
		entity.ServiceId,
		entity.Date,
		entity.Hour,
		entity.CreatedAt,
		entity.Ttl,
	)

	if err != nil {
		return types.ApiError{
			Msg:      err.Error(),
			Function: "Save -> r.connection.Exec()",
			File:     "pg-booking-session-repository.go",
			Values: []string{
				query.String(),
				entity.Id,
				entity.BusinessId,
				strconv.Itoa(entity.ChatId),
			},
		}
	}

	return nil
}

func (r *PgBookingSessionRepository) Update(entity types.BookingSession) error {
	var query strings.Builder

	query.WriteString(`UPDATE booking_session `)
	query.WriteString(`SET business_id = $2, chat_id = $3, service_id = $4, date = $5, hour = $6, `)
	query.WriteString(`created_at = $7, ttl = $8 WHERE id = $1;`)

	_, err := r.connection.Exec(
		query.String(),
		entity.Id,
		entity.BusinessId,
		entity.ChatId,
		entity.ServiceId,
		entity.Date,
		entity.Hour,
		entity.CreatedAt,
		entity.Ttl,
	)

	if err != nil {
		return types.ApiError{
			Msg:      err.Error(),
			Function: "Update -> r.connection.Exec()",
			File:     "pg-booking-session-repository.go",
			Values: []string{
				query.String(),
				entity.Id,
				entity.BusinessId,
				strconv.Itoa(entity.ChatId),
			},
		}
	}

	return nil
}
