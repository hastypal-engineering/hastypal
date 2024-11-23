package repository

import (
	"database/sql"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	types "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
)

type PgGoogleTokenRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgGoogleTokenRepository(connection *sql.DB) *PgGoogleTokenRepository {
	transformer := helper.NewCriteriaToSqlService("google_token")

	return &PgGoogleTokenRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgGoogleTokenRepository) Find(criteria types.Criteria) ([]types.GoogleToken, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, types.ApiError{
			Msg:      err.Error(),
			Function: "Find -> r.transformer.Transform()",
			File:     "pg-google-token-repository.go",
		}
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, types.ApiError{
			Msg:      queryErr.Error(),
			Function: "Find -> r.connection.Query()",
			File:     "pg-google-token-repository.go",
			Values:   []string{query},
		}
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
			return nil, types.ApiError{
				Msg:      scanErr.Error(),
				Function: "Find -> rows.Scan()",
				File:     "pg-google-token-repository.go",
			}
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
		return types.GoogleToken{}, types.ApiError{
			Msg:      err.Error(),
			Function: "FindOne -> r.transformer.Transform()",
			File:     "pg-google-token-repository.go",
		}
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
		if errors.As(err, &sql.ErrNoRows) {
			return types.GoogleToken{}, types.ApiError{
				Msg:      "Entity GoogleToken not found",
				Function: "FindOne -> rows.Scan()",
				File:     "pg-google-token-repository.go",
				Values:   []string{query},
				Domain:   true,
			}
		}

		return types.GoogleToken{}, types.ApiError{
			Msg:      scanErr.Error(),
			Function: "FindOne -> rows.Scan()",
			File:     "pg-google-token-repository.go",
			Values:   []string{query},
		}
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
		return types.ApiError{
			Msg:      err.Error(),
			Function: "Save -> r.connection.Exec()",
			File:     "pg-google-token-repository.go",
			Values: []string{
				query.String(),
				entity.BusinessId,
				entity.CreatedAt,
				entity.UpdatedAt,
			},
		}
	}

	return nil
}

func (r *PgGoogleTokenRepository) Update(_ types.GoogleToken) error {
	return types.ApiError{
		Msg:      "Method not implemented yet",
		Function: "Update",
		File:     "pg-google-token-repository.go",
	}
}
