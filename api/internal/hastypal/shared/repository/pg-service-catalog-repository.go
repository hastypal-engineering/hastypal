package repository

import (
	"database/sql"
	"strings"

	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
)

type PgServiceCatalogRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgServiceCatalogRepository(connection *sql.DB) *PgServiceCatalogRepository {
	transformer, _ := helper.NewCriteriaToSqlService((&types.ServiceCatalog{}))

	return &PgServiceCatalogRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgServiceCatalogRepository) Find(criteria types.Criteria) error {
	return exception.
		New("Method not implemented").
		Trace("Find", "pg-service-catalog-repository.go")
}

func (r *PgServiceCatalogRepository) FindOne(criteria types.Criteria) error {
	return exception.
		New("Method not implemented").
		Trace("FindOne", "pg-service-catalog-repository.go")
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
		return exception.
			New(err.Error()).
			Trace("r.connection.Exec", "pg-service-catalog-repository.go").
			WithValues([]string{query.String(), entity.Id, entity.Name, entity.BusinessId})
	}

	return nil
}

func (r *PgServiceCatalogRepository) Update(entity types.ServiceCatalog) error {
	return exception.
		New("Method not implemented").
		Trace("Update", "pg-service-catalog-repository.go")
}

func (r *PgServiceCatalogRepository) Delete(criteria types.Criteria) error {
	return exception.
		New("Method not implemented").
		Trace("Delete", "pg-service-catalog-repository.go")
}
