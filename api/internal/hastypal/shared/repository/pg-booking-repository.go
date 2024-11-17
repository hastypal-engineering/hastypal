package repository

import (
	"database/sql"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
)

type PgBookingRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgBookingRepository(connection *sql.DB) *PgBookingRepository {
	transformer := helper.NewCriteriaToSqlService("booking")

	return &PgBookingRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgBookingRepository) Find(criteria types2.Criteria) ([]types2.Booking, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, types2.ApiError{
			Msg:      err.Error(),
			Function: "Find -> r.transformer.Transform()",
			File:     "pg-booking-repository.go",
		}
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, types2.ApiError{
			Msg:      queryErr.Error(),
			Function: "Find -> r.connection.Query()",
			File:     "pg-booking-repository.go",
			Values:   []string{query},
		}
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

	var results []types2.Booking

	for rows.Next() {
		if scanErr := rows.Scan(
			&id,
			&session_id,
			&business_id,
			&service_id,
			&booking_date,
			&created_at,
		); scanErr != nil {
			return nil, types2.ApiError{
				Msg:      scanErr.Error(),
				Function: "Find -> rows.Scan()",
				File:     "pg-booking-repository.go",
			}
		}

		results = append(results, types2.Booking{
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

func (r *PgBookingRepository) FindOne(criteria types2.Criteria) (types2.Booking, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types2.Booking{}, types2.ApiError{
			Msg:      err.Error(),
			Function: "FindOne -> r.transformer.Transform()",
			File:     "pg-booking-repository.go",
		}
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
		if errors.As(err, &sql.ErrNoRows) {
			return types2.Booking{}, types2.ApiError{
				Msg:      "Entity Business not found",
				Function: "FindOne -> rows.Scan()",
				File:     "pg-booking-repository.go",
				Values:   []string{query},
				Domain:   true,
			}
		}

		return types2.Booking{}, types2.ApiError{
			Msg:      scanErr.Error(),
			Function: "FindOne -> rows.Scan()",
			File:     "pg-booking-repository.go",
			Values:   []string{query},
		}
	}

	return types2.Booking{
		Id:         id,
		SessionId:  session_id,
		BusinessId: business_id,
		ServiceId:  service_id,
		When:       booking_date,
		CreatedAt:  created_at,
	}, nil
}

func (r *PgBookingRepository) Save(entity types2.Booking) error {
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
		return types2.ApiError{
			Msg:      err.Error(),
			Function: "Save -> r.connection.Exec()",
			File:     "pg-booking-repository.go",
			Values: []string{
				query.String(),
				entity.Id,
				entity.BusinessId,
				entity.SessionId,
			},
		}
	}

	return nil
}

func (r *PgBookingRepository) Update(_ types2.Booking) error {
	return types2.ApiError{
		Msg:      "Method not implemented yet",
		Function: "Update",
		File:     "pg-booking-repository.go",
	}
}