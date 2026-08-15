package google

import (
	"context"
	"database/sql"
)

type GoogleRepository interface {
	Save(ctx context.Context, token *GoogleToken) error
	GetByBusinessID(ctx context.Context, ID int) (*GoogleToken, error)
}

type PgGoogleRepository struct {
	connection *sql.DB
}

func NewPgGoogleRepository(connection *sql.DB) *PgGoogleRepository {
	return &PgGoogleRepository{
		connection: connection,
	}
}

func (r *PgGoogleRepository) Save(ctx context.Context, token *GoogleToken) error {
	return nil
}

func (r *PgGoogleRepository) GetByBusinessID(ctx context.Context, ID int) (*GoogleToken, error) {
	return nil, nil
}
