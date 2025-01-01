package repository

import (
	"database/sql"

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
	return exception.
		New("Method not implemented").
		Trace("Save", "pg-service-catalog-repository.go")
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
