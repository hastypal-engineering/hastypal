package repository

import (
	"database/sql"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
)

type PgGoogleTokenRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgGoogleTokenRepository(connection *sql.DB) *PgGoogleTokenRepository {
	transformer, _ := helper.NewCriteriaToSqlService(&types.GoogleToken{})

	return &PgGoogleTokenRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgGoogleTokenRepository) Find(criteria types.Criteria) ([]types.GoogleToken, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, exception.
			New(err.Error()).
			Trace("r.transformer.Transform", "pg-google-token-repository.go")
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, exception.
			New(queryErr.Error()).
			Trace("r.connection.Query", "pg-google-token-repository.go").
			WithValues([]string{query})
	}

	defer rows.Close()

	var (
		business_id   string
		access_token  string
		token_type    string
		refresh_token string
		created_at    string
		updated_at    string
	)

	var results []types.GoogleToken

	for rows.Next() {
		if scanErr := rows.Scan(
			&business_id,
			&access_token,
			&token_type,
			&refresh_token,
			&created_at,
			&updated_at,
		); scanErr != nil {
			return nil, exception.
				New(scanErr.Error()).
				Trace("rows.Scan", "pg-google-token-repository.go").
				WithValues([]string{query})
		}

		results = append(results, types.GoogleToken{
			BusinessId:   business_id,
			AccessToken:  access_token,
			TokenType:    token_type,
			RefreshToken: refresh_token,
			CreatedAt:    created_at,
			UpdatedAt:    updated_at,
		})
	}

	return results, nil
}

func (r *PgGoogleTokenRepository) FindOne(criteria types.Criteria) (types.GoogleToken, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types.GoogleToken{}, exception.
			New(err.Error()).
			Trace("r.transformer.Transform", "pg-google-token-repository.go")
	}

	var (
		business_id   string
		access_token  string
		token_type    string
		refresh_token string
		created_at    string
		updated_at    string
	)

	if scanErr := r.connection.QueryRow(query).Scan(
		&business_id,
		&access_token,
		&token_type,
		&refresh_token,
		&created_at,
		&updated_at,
	); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return types.GoogleToken{}, exception.
				New("Google token not found").
				Trace("r.connection.QueryRow", "pg-google-token-repository.go").
				WithValues([]string{query}).
				Domain()
		}

		return types.GoogleToken{}, exception.
			New(scanErr.Error()).
			Trace("r.connection.QueryRow", "pg-google-token-repository.go").
			WithValues([]string{query})
	}

	return types.GoogleToken{
		BusinessId:   business_id,
		AccessToken:  access_token,
		TokenType:    token_type,
		RefreshToken: refresh_token,
		CreatedAt:    created_at,
		UpdatedAt:    updated_at,
	}, nil
}

func (r *PgGoogleTokenRepository) Save(entity types.GoogleToken) error {
	var query strings.Builder

	query.WriteString(`INSERT INTO google_token `)
	query.WriteString(`(business_id, access_token, token_type, refresh_token, created_at, updated_at) `)
	query.WriteString(`VALUES ($1, $2, $3, $4, $5, $6);`)

	_, err := r.connection.Exec(
		query.String(),
		entity.BusinessId,
		entity.AccessToken,
		entity.TokenType,
		entity.RefreshToken,
		entity.CreatedAt,
		entity.UpdatedAt,
	)

	if err != nil {
		return exception.
			New(err.Error()).
			Trace("r.connection.Exec", "pg-google-token-repository.go").
			WithValues([]string{
				query.String(),
				entity.BusinessId,
				entity.CreatedAt,
				entity.UpdatedAt,
			})
	}

	return nil
}

func (r *PgGoogleTokenRepository) Update(_ types.GoogleToken) error {
	return exception.
		New("Method not implemented").
		Trace("Update", "pg-google-token-repository.go")
}
