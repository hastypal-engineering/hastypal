package repository

import (
	"database/sql"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	types "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
)

type PgServiceCatalogRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgServiceCatalogRepository(connection *sql.DB) *PgServiceCatalogRepository {
	transformer := helper.NewCriteriaToSqlService("service_catalog")

	return &PgServiceCatalogRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgServiceCatalogRepository) Save(entity types.ServiceCatalog) error {
	var query strings.Builder

	query.WriteString(`INSERT INTO service_catalog `)
	query.WriteString(`(id, name, price, currency, duration, business_id) `)
	query.WriteString(`VALUES ($1, $2, $3, $4, $5, $6);`)

	_, err := r.connection.Exec(
		query.String(),
		entity.Id,
		entity.Name,
		entity.Price,
		entity.Currency,
		entity.Duration,
		entity.BusinessId,
	)

	if err != nil {
		return types.ApiError{
			Msg:      err.Error(),
			Function: "Save -> r.connection.Exec()",
			File:     "pg-service-catalog-repository.go",
			Values: []string{
				query.String(),
				entity.Id,
				entity.Name,
				entity.BusinessId,
			},
		}
	}

	return nil
}
