package repository

import (
	"database/sql"
	"errors"

	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
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
		id              string
		name            string
		contact_phone   string
		email           string
		password        string
		service_catalog []types.ServiceCatalog
		opening_hours   map[string][]string
		channel_name    string
		location        string
		created_at      string
		updated_at      string
	)

	var results []types.Business

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
			return nil, types.ApiError{
				Msg:      scanErr.Error(),
				Function: "Find -> rows.Scan()",
				File:     "pg-business-repository.go",
			}
		}

		results = append(results, types.Business{
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
		id                  string
		name                string
		communication_phone string
		email               string
		password            string
		service_catalog     []types.ServiceCatalog
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
			return types.Business{}, types.ApiError{
				Msg:      "Entity Business not found",
				Function: "FindOne -> rows.Scan()",
				File:     "pg-business-repository.go",
				Values:   []string{query},
				Domain:   true,
			}
		}

		return types.Business{}, types.ApiError{
			Msg:      scanErr.Error(),
			Function: "FindOne -> rows.Scan()",
			File:     "pg-business-repository.go",
			Values:   []string{query},
		}
	}

	return types.Business{
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

func (r *PgBusinessRepository) Save(entity types.Business) error {
	var query = `INSERT INTO business (id, name, communication_phone, email, password, service_catalog, opening_hours, holidays, channel_name, location, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.connection.Exec(
		query,
		entity.Id,
		entity.Name,
		entity.ContactPhone,
		entity.Email,
		entity.Password,
		entity.ServiceCatalog,
		entity.OpeningHours,
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
				query,
				entity.Id,
				entity.Name,
				entity.ContactPhone,
				entity.Email,
			},
		}
	}

	return nil
}
