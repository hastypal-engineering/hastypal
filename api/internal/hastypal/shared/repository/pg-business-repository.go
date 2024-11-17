package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
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

func (r *PgBusinessRepository) Find(criteria types2.Criteria) ([]types2.Business, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, types2.ApiError{
			Msg:      err.Error(),
			Function: "Find -> r.transformer.Transform()",
			File:     "pg-business-repository.go",
		}
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, types2.ApiError{
			Msg:      queryErr.Error(),
			Function: "Find -> r.connection.Query()",
			File:     "pg-business-repository.go",
			Values:   []string{query},
		}
	}

	defer rows.Close()

	var (
		id              string
		name            string
		contact_phone   string
		email           string
		password        string
		service_catalog []types2.ServiceCatalog
		opening_hours   map[string][]string
		channel_name    string
		location        string
		created_at      string
		updated_at      string
	)

	var results []types2.Business

	for rows.Next() {
		if scanErr := rows.Scan(
			&id,
			&name,
			&contact_phone,
			&email,
			&password,
			&service_catalog,
			&opening_hours,
			&channel_name,
			&location,
			&created_at,
			&updated_at,
		); scanErr != nil {
			return nil, types2.ApiError{
				Msg:      scanErr.Error(),
				Function: "Find -> rows.Scan()",
				File:     "pg-business-repository.go",
			}
		}

		results = append(results, types2.Business{
			Id:             id,
			Name:           name,
			ContactPhone:   contact_phone,
			Email:          email,
			Password:       password,
			ServiceCatalog: service_catalog,
			OpeningHours:   opening_hours,
			ChannelName:    channel_name,
			Location:       location,
			CreatedAt:      created_at,
			UpdatedAt:      updated_at,
		})
	}

	return results, nil
}

func (r *PgBusinessRepository) FindOne(criteria types2.Criteria) (types2.Business, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types2.Business{}, types2.ApiError{
			Msg:      err.Error(),
			Function: "FindOne -> r.transformer.Transform()",
			File:     "pg-business-repository.go",
		}
	}

	var (
		id                  string
		name                string
		communication_phone string
		email               string
		password            string
		service_catalog     []types2.ServiceCatalog
		opening_hours       map[string][]string
		channel_name        string
		location            string
		created_at          string
		updated_at          string
	)

	if scanErr := r.connection.QueryRow(query).Scan(
		&id,
		&name,
		&communication_phone,
		&email,
		&password,
		&service_catalog,
		&opening_hours,
		&channel_name,
		&location,
		&created_at,
		&updated_at,
	); scanErr != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return types2.Business{}, types2.ApiError{
				Msg:      "Entity Business not found",
				Function: "FindOne -> rows.Scan()",
				File:     "pg-business-repository.go",
				Values:   []string{query},
				Domain:   true,
			}
		}

		return types2.Business{}, types2.ApiError{
			Msg:      scanErr.Error(),
			Function: "FindOne -> rows.Scan()",
			File:     "pg-business-repository.go",
			Values:   []string{query},
		}
	}

	return types2.Business{
		Id:             id,
		Name:           name,
		ContactPhone:   communication_phone,
		Email:          email,
		Password:       password,
		ServiceCatalog: service_catalog,
		OpeningHours:   opening_hours,
		ChannelName:    channel_name,
		Location:       location,
		CreatedAt:      created_at,
		UpdatedAt:      updated_at,
	}, nil
}

func (r *PgBusinessRepository) Save(entity types2.Business) error {
	var query strings.Builder

	query.WriteString(`INSERT INTO business `)
	query.WriteString(`(id, name, contact_phone, email, password, opening_hours, holidays, channel_name, location, created_at, updated_at) `)
	query.WriteString(`VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`)

	openingHours, openingHoursErr := json.Marshal(entity.OpeningHours)
	if openingHoursErr != nil {
		return types2.ApiError{
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
		return types2.ApiError{
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
		return types2.ApiError{
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

func (r *PgBusinessRepository) Update(_ types2.Business) error {
	return types2.ApiError{
		Msg:      "Method not implemented yet",
		Function: "Update",
		File:     "pg-business-repository.go",
	}
}
