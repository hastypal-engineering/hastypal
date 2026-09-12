package business

import (
	"context"
	"database/sql"
	"time"

	"github.com/adriein/hastypal/database"
	"github.com/rotisserie/eris"
)

var BusinessNotFound = eris.New("Business not found")

type BusinessRepository interface {
	Create(ctx context.Context, business *Business) (int, error)
	CreateService(ctx context.Context, service *ServiceCatalog) (int, error)
	GetByID(ctx context.Context, ID int) (*Business, error)
	GetByPublicID(ctx context.Context, ID string) (*Business, error)
}

type PgBusinessRepository struct {
	connection *sql.DB
}

func NewPgBusinessRepository(connection *sql.DB) *PgBusinessRepository {
	return &PgBusinessRepository{
		connection: connection,
	}
}

func (r *PgBusinessRepository) Create(ctx context.Context, business *Business) (int, error) {
	query := `
		INSERT INTO ha_business (
			hab_public_id,
			hab_name,
			hab_contact_phone,
			hab_email,
			hab_address,
			hab_country,
			hab_lang,
			hab_date_add,
			hab_date_upd
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING hab_id;
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	now := time.Now()

	err := r.connection.QueryRowContext(ctxTimeout, query,
		business.PublicID,
		business.Name,
		business.ContactPhone,
		business.Email,
		business.Address,
		business.Country,
		business.Lang,
		now,
		now,
	).Scan(&business.ID)
	if err != nil {
		return 0, eris.Wrap(err, "Failed to create business")
	}

	return business.ID, nil
}

func (r *PgBusinessRepository) CreateService(ctx context.Context, service *ServiceCatalog) (int, error) {
	query := `
	INSERT INTO ha_service_catalog (
			hasc_name,
			hasc_description,
			hasc_price,
			hasc_currency,
			hasc_duration,
			hasc_business_id,
			hasc_date_add,
			hasc_date_upd
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING hasc_id;
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	now := time.Now()

	err := r.connection.QueryRowContext(ctxTimeout, query,
		service.Name,
		service.Description,
		service.Price,
		service.Currency,
		service.Duration,
		service.BusinessID,
		now,
		now,
	).Scan(&service.ID)
	if err != nil {
		return 0, eris.Wrap(err, "Failed to create service")
	}

	return service.ID, nil
}

func (r *PgBusinessRepository) GetByID(ctx context.Context, ID int) (*Business, error) {
	query := `
		SELECT
			b.hab_id,
			b.hab_public_id,
			b.hab_name,
			b.hab_contact_phone,
			b.hab_email,
			b.hab_address,
			b.hab_country,
			b.hab_lang,
			b.hab_date_add,
			b.hab_date_upd,
			s.hasc_id,
			s.hasc_name,
			s.hasc_description,
			s.hasc_price,
			s.hasc_currency,
			s.hasc_duration,
			s.hasc_business_id,
			s.hasc_date_add,
			s.hasc_date_upd
		FROM
			ha_business b
		INNER JOIN
			ha_service_catalog s ON s.hasc_business_id = b.hab_id
		WHERE
			b.hab_id = $1;
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rows, err := r.connection.QueryContext(ctxTimeout, query, ID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to query business by ID")
	}

	defer database.CloseRowsSafely(rows, &err)

	business := &Business{}

	for rows.Next() {
		var service ServiceCatalog

		err := rows.Scan(
			&business.ID,
			&business.PublicID,
			&business.Name,
			&business.ContactPhone,
			&business.Email,
			&business.Address,
			&business.Country,
			&business.Lang,
			&business.DateAdd,
			&business.DateUpd,
			&service.ID,
			&service.Name,
			&service.Description,
			&service.Price,
			&service.Currency,
			&service.Duration,
			&service.BusinessID,
			&service.DateAdd,
			&service.DateUpd,
		)
		if err != nil {
			return nil, eris.Wrap(err, "Failed to scan business by ID")
		}

		business.ServiceCatalog = append(business.ServiceCatalog, &service)
	}

	if err := rows.Err(); err != nil {
		return nil, eris.Wrap(err, "Failed to iterate business rows")
	}

	if business.ID == 0 {
		return nil, BusinessNotFound
	}

	return business, nil
}

func (r *PgBusinessRepository) GetByPublicID(ctx context.Context, ID string) (*Business, error) {
	query := `
		SELECT
			b.hab_id,
			b.hab_public_id,
			b.hab_name,
			b.hab_contact_phone,
			b.hab_email,
			b.hab_address,
			b.hab_country,
			b.hab_lang,
			b.hab_date_add,
			b.hab_date_upd,
			s.hasc_id,
			s.hasc_name,
			s.hasc_description,
			s.hasc_price,
			s.hasc_currency,
			s.hasc_duration,
			s.hasc_business_id,
			s.hasc_date_add,
			s.hasc_date_upd
		FROM
			ha_business b
		INNER JOIN
			ha_service_catalog s ON s.hasc_business_id = b.hab_id
		WHERE
			b.hab_public_id = $1;
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rows, err := r.connection.QueryContext(ctxTimeout, query, ID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to query business by public ID")
	}

	defer database.CloseRowsSafely(rows, &err)

	business := &Business{}

	for rows.Next() {
		var service ServiceCatalog

		err := rows.Scan(
			&business.ID,
			&business.PublicID,
			&business.Name,
			&business.ContactPhone,
			&business.Email,
			&business.Address,
			&business.Country,
			&business.Lang,
			&business.DateAdd,
			&business.DateUpd,
			&service.ID,
			&service.Name,
			&service.Description,
			&service.Price,
			&service.Currency,
			&service.Duration,
			&service.BusinessID,
			&service.DateAdd,
			&service.DateUpd,
		)
		if err != nil {
			return nil, eris.Wrap(err, "Failed to scan business by public ID")
		}

		business.ServiceCatalog = append(business.ServiceCatalog, &service)
	}

	if err := rows.Err(); err != nil {
		return nil, eris.Wrap(err, "Failed to iterate business rows")
	}

	if business.ID == 0 {
		return nil, BusinessNotFound
	}

	return business, nil
}
