package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
)

type PgBusinessRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgBusinessRepository(connection *sql.DB) *PgBusinessRepository {
	transformer := helper.NewCriteriaToSqlService("business")

	return &PgBusinessRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgBusinessRepository) Find(criteria types.Criteria) ([]types.Business, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, types.ApiError{
			Msg:      err.Error(),
			Function: "Find -> r.transformer.Transform()",
			File:     "pg-business-repository.go",
		}
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, types.ApiError{
			Msg:      queryErr.Error(),
			Function: "Find -> r.connection.Query()",
			File:     "pg-business-repository.go",
			Values:   []string{query},
		}
	}

	defer rows.Close()

	var (
		id            string
		name          string
		contact_phone string
		email         string
		password      string
		channel_name  string
		location      string
		opening_hours []uint8
		holidays      []uint8
		created_at    string
		updated_at    string
	)

	var results []types.Business

	for rows.Next() {
		if scanErr := rows.Scan(
			&id,
			&name,
			&contact_phone,
			&email,
			&password,
			&channel_name,
			&location,
			&opening_hours,
			&holidays,
			&created_at,
			&updated_at,
		); scanErr != nil {
			return nil, types.ApiError{
				Msg:      scanErr.Error(),
				Function: "Find -> rows.Scan()",
				File:     "pg-business-repository.go",
			}
		}

		var openingHours map[string][]string

		openingHoursUnMarshalErr := json.Unmarshal(opening_hours, &openingHours)

		if openingHoursUnMarshalErr != nil {
			return results, types.ApiError{
				Msg:      openingHoursUnMarshalErr.Error(),
				Function: "FindOne -> json.Unmarshal(openingHours)",
				File:     "pg-business-repository.go",
				Values:   []string{query},
			}
		}

		var holidaysMap map[string][]string

		holidaysUnMarshalErr := json.Unmarshal(holidays, &holidaysMap)

		if holidaysUnMarshalErr != nil {
			return results, types.ApiError{
				Msg:      holidaysUnMarshalErr.Error(),
				Function: "FindOne -> json.Unmarshal(holidays)",
				File:     "pg-business-repository.go",
				Values:   []string{query},
			}
		}

		results = append(results, types.Business{
			Id:           id,
			Name:         name,
			ContactPhone: contact_phone,
			Email:        email,
			Password:     password,
			OpeningHours: openingHours,
			Holidays:     holidaysMap,
			ChannelName:  channel_name,
			Location:     location,
			CreatedAt:    created_at,
			UpdatedAt:    updated_at,
		})
	}

	return results, nil
}

func (r *PgBusinessRepository) FindOne(criteria types.Criteria) (types.Business, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types.Business{}, types.ApiError{
			Msg:      err.Error(),
			Function: "FindOne -> r.transformer.Transform()",
			File:     "pg-business-repository.go",
		}
	}

	var (
		id            string
		name          string
		contact_phone string
		email         string
		password      string
		channel_name  string
		location      string
		opening_hours []uint8
		holidays      []uint8
		created_at    string
		updated_at    string
	)

	if scanErr := r.connection.QueryRow(query).Scan(
		&id,
		&name,
		&contact_phone,
		&email,
		&password,
		&channel_name,
		&location,
		&opening_hours,
		&holidays,
		&created_at,
		&updated_at,
	); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return types.Business{}, types.ApiError{
				Msg:      "Entity Business not found",
				Function: "r.connection.QueryRow.Scan",
				File:     "pg-business-repository.go",
				Values:   []string{query},
				Domain:   true,
			}
		}

		return types.Business{}, exception.
			New("Entity Business not found").
			Trace("r.connection.QueryRow.Scan", "pg-business-repository.go").
			WithValues([]string{query}).
			Domain()
	}

	var openingHours map[string][]string

	openingHoursUnMarshalErr := json.Unmarshal(opening_hours, &openingHours)

	if openingHoursUnMarshalErr != nil {
		return types.Business{}, types.ApiError{
			Msg:      openingHoursUnMarshalErr.Error(),
			Function: "FindOne -> json.Unmarshal(openingHours)",
			File:     "pg-business-repository.go",
			Values:   []string{query},
		}
	}

	var holidaysMap map[string][]string

	holidaysUnMarshalErr := json.Unmarshal(holidays, &holidaysMap)

	if holidaysUnMarshalErr != nil {
		return types.Business{}, types.ApiError{
			Msg:      holidaysUnMarshalErr.Error(),
			Function: "FindOne -> json.Unmarshal(holidays)",
			File:     "pg-business-repository.go",
			Values:   []string{query},
		}
	}

	return types.Business{
		Id:           id,
		Name:         name,
		ContactPhone: contact_phone,
		Email:        email,
		Password:     password,
		OpeningHours: openingHours,
		Holidays:     holidaysMap,
		ChannelName:  channel_name,
		Location:     location,
		CreatedAt:    created_at,
		UpdatedAt:    updated_at,
	}, nil
}

func (r *PgBusinessRepository) Save(entity types.Business) error {
	var query strings.Builder

	query.WriteString(`INSERT INTO business `)
	query.WriteString(`(id, name, contact_phone, email, password, opening_hours, holidays, channel_name, location, created_at, updated_at) `)
	query.WriteString(`VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`)

	openingHours, openingHoursErr := json.Marshal(entity.OpeningHours)
	if openingHoursErr != nil {
		return types.ApiError{
			Msg:      openingHoursErr.Error(),
			Function: "Save -> json.Marshal(entity.OpeningHours)",
			File:     "pg-business-repository.go",
			Values: []string{
				query.String(),
				entity.Id,
				entity.Name,
				entity.ContactPhone,
				entity.Email,
			},
		}
	}

	holidays, holidaysErr := json.Marshal(entity.Holidays)
	if holidaysErr != nil {
		return types.ApiError{
			Msg:      holidaysErr.Error(),
			Function: "Save -> json.Marshal(entity.Holidays)",
			File:     "pg-business-repository.go",
			Values: []string{
				query.String(),
				entity.Id,
				entity.Name,
				entity.ContactPhone,
				entity.Email,
			},
		}
	}

	_, err := r.connection.Exec(
		query.String(),
		entity.Id,
		entity.Name,
		entity.ContactPhone,
		entity.Email,
		entity.Password,
		openingHours,
		holidays,
		entity.ChannelName,
		entity.Location,
		entity.CreatedAt,
		entity.UpdatedAt,
	)

	if err != nil {
		return types.ApiError{
			Msg:      err.Error(),
			Function: "Save -> r.connection.Exec()",
			File:     "pg-business-repository.go",
			Values: []string{
				query.String(),
				entity.Id,
				entity.Name,
				entity.ContactPhone,
				entity.Email,
			},
		}
	}

	return nil
}

func (r *PgBusinessRepository) Update(_ types.Business) error {
	return types.ApiError{
		Msg:      "Method not implemented yet",
		Function: "Update",
		File:     "pg-business-repository.go",
	}
}
