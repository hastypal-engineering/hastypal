package repository

import (
	"database/sql"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
)

type PgBookingRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgBookingRepository(connection *sql.DB) *PgBookingRepository {
	transformer, _ := helper.NewCriteriaToSqlService(&types.Booking{})

	return &PgBookingRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgBookingRepository) Find(criteria types.Criteria) ([]types.Booking, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, exception.
			New(err.Error()).
			Trace("r.transformer.Transform", "pg-booking-repository.go")
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, exception.
			New(queryErr.Error()).
			Trace("r.connection.Query", "pg-booking-repository.go")
	}

	defer rows.Close()

	var (
		id           string
		session_id   string
		business_id  string
		service_id   string
		booking_date string
		created_at   string
	)

	var results []types.Booking

	for rows.Next() {
		if scanErr := rows.Scan(
			&id,
			&session_id,
			&business_id,
			&service_id,
			&booking_date,
			&created_at,
		); scanErr != nil {
			return nil, exception.
				New(scanErr.Error()).
				Trace("rows.Scan", "pg-booking-repository.go")
		}

		results = append(results, types.Booking{
			Id:         id,
			SessionId:  session_id,
			BusinessId: business_id,
			ServiceId:  service_id,
			When:       booking_date,
			CreatedAt:  created_at,
		})
	}

	return results, nil
}

func (r *PgBookingRepository) FindOne(criteria types.Criteria) (types.Booking, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types.Booking{}, exception.
			New(err.Error()).
			Trace("r.transformer.Transform", "pg-booking-repository.go")
	}

	var (
		id           string
		session_id   string
		business_id  string
		service_id   string
		booking_date string
		created_at   string
	)

	if scanErr := r.connection.QueryRow(query).Scan(
		&id,
		&session_id,
		&business_id,
		&service_id,
		&booking_date,
		&created_at,
	); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return types.Booking{}, exception.
				New("Booking not found").
				Trace("r.connection.QueryRow", "pg-booking-repository.go").
				WithValues([]string{query}).
				Domain()
		}

		return types.Booking{}, exception.
			New(scanErr.Error()).
			Trace("r.connection.QueryRow", "pg-booking-repository.go").
			WithValues([]string{query})
	}

	return types.Booking{
		Id:         id,
		SessionId:  session_id,
		BusinessId: business_id,
		ServiceId:  service_id,
		When:       booking_date,
		CreatedAt:  created_at,
	}, nil
}

func (r *PgBookingRepository) Save(entity types.Booking) error {
	var query strings.Builder

	query.WriteString(`INSERT INTO booking `)
	query.WriteString(`(id, session_id, business_id, service_id, booking_date, created_at) `)
	query.WriteString(`VALUES ($1, $2, $3, $4, $5, $6);`)

	_, err := r.connection.Exec(
		query.String(),
		entity.Id,
		entity.SessionId,
		entity.BusinessId,
		entity.ServiceId,
		entity.When,
		entity.CreatedAt,
	)

	if err != nil {
		return exception.
			New(err.Error()).
			Trace("r.connection.Exec", "pg-booking-repository.go").
			WithValues([]string{
				query.String(),
				entity.Id,
				entity.BusinessId,
				entity.SessionId,
			})
	}

	return nil
}

func (r *PgBookingRepository) Update(_ types.Booking) error {
	return exception.
		New("Method not implemented").
		Trace("Update", "pg-booking-repository.go")
}
