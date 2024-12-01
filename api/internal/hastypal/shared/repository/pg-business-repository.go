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
		return nil, exception.New(err.Error()).
			Trace("r.transformer.Transform", "pg-business-repository.go")
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, exception.New(queryErr.Error()).
			Trace("r.connection.Query", "pg-business-repository.go").
			WithValues([]string{query})
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
			return nil, exception.New(scanErr.Error()).
				Trace("rows.Scan", "pg-business-repository.go").
				WithValues([]string{query})
		}

		var openingHours map[string][]string

		openingHoursUnMarshalErr := json.Unmarshal(opening_hours, &openingHours)

		if openingHoursUnMarshalErr != nil {
			return results, exception.New(openingHoursUnMarshalErr.Error()).
				Trace("json.Unmarshal(openingHours)", "pg-business-repository.go")
		}

		var holidaysMap map[string][]string

		holidaysUnMarshalErr := json.Unmarshal(holidays, &holidaysMap)

		if holidaysUnMarshalErr != nil {
			return results, exception.New(holidaysUnMarshalErr.Error()).
				Trace("json.Unmarshal(holidays)", "pg-business-repository.go")
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
		return types.Business{}, exception.New(err.Error()).
			Trace("r.transformer.Transform", "pg-business-repository.go")
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
			return types.Business{}, exception.
				New("Entity Business not found").
				Trace("r.connection.QueryRow.Scan", "pg-business-repository.go").
				WithValues([]string{query}).
				Domain()
		}

		return types.Business{}, exception.
			New(scanErr.Error()).
			Trace("r.connection.QueryRow.Scan", "pg-business-repository.go").
			WithValues([]string{query})
	}

	var openingHours map[string][]string

	openingHoursUnMarshalErr := json.Unmarshal(opening_hours, &openingHours)

	if openingHoursUnMarshalErr != nil {
		return types.Business{}, exception.
			New(openingHoursUnMarshalErr.Error()).
			Trace("json.Unmarshal(openingHours)", "pg-business-repository.go").
			WithValues([]string{query})
	}

	var holidaysMap map[string][]string

	holidaysUnMarshalErr := json.Unmarshal(holidays, &holidaysMap)

	if holidaysUnMarshalErr != nil {
		return types.Business{}, exception.
			New(holidaysUnMarshalErr.Error()).
			Trace("json.Unmarshal(holidays, &holidaysMap)", "pg-business-repository.go").
			WithValues([]string{query})

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
		return exception.New(openingHoursErr.Error()).
			Trace("json.Marshal(entity.OpeningHours)", "pg-business-repository.go").
			WithValues([]string{query.String(), entity.Id, entity.Name, entity.ContactPhone, entity.Email})
	}

	holidays, holidaysErr := json.Marshal(entity.Holidays)

	if holidaysErr != nil {
		return exception.New(holidaysErr.Error()).
			Trace("json.Marshal(entity.Holidays)", "pg-business-repository.go").
			WithValues([]string{query.String(), entity.Id, entity.Name, entity.ContactPhone, entity.Email})
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
		return exception.New(err.Error()).
			Trace("r.connection.Exec", "pg-business-repository.go").
			WithValues([]string{query.String(), entity.Id, entity.Name, entity.ContactPhone, entity.Email})
	}

	return nil
}

func (r *PgBusinessRepository) Update(_ types.Business) error {
	return exception.New("Method not implemented").
		Trace("Update", "pg-business-repository.go")
}
