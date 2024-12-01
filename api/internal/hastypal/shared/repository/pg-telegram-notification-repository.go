package repository

import (
	"database/sql"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"strings"
)

type PgTelegramNotificationRepository struct {
	connection  *sql.DB
	transformer *helper.CriteriaToSqlService
}

func NewPgTelegramNotificationRepository(connection *sql.DB) *PgTelegramNotificationRepository {
	transformer := helper.NewCriteriaToSqlService("telegram_notification")

	return &PgTelegramNotificationRepository{
		connection:  connection,
		transformer: transformer,
	}
}

func (r *PgTelegramNotificationRepository) Find(criteria types.Criteria) ([]types.TelegramNotification, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, exception.
			New(err.Error()).
			Trace("r.transformer.Transform", "pg-telegram-notification-repository.go")
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, exception.
			New(queryErr.Error()).
			Trace("r.connection.Query", "pg-telegram-notification-repository.go").
			WithValues([]string{query})
	}

	defer rows.Close()

	var (
		id            string
		session_id    string
		business_id   string
		booking_id    string
		scheduled_at  string
		chat_id       int
		business_name string
		created_at    string
	)

	var results []types.TelegramNotification

	for rows.Next() {
		if scanErr := rows.Scan(
			&id,
			&session_id,
			&business_id,
			&booking_id,
			&scheduled_at,
			&chat_id,
			&business_name,
			&created_at,
		); scanErr != nil {
			return nil, exception.
				New(scanErr.Error()).
				Trace("rows.Scan", "pg-telegram-notification-repository.go").
				WithValues([]string{query})
		}

		results = append(results, types.TelegramNotification{
			Id:          id,
			SessionId:   session_id,
			BusinessId:  business_id,
			BookingId:   booking_id,
			ScheduledAt: scheduled_at,
			ChatId:      chat_id,
			From:        business_name,
			CreatedAt:   created_at,
		})
	}

	return results, nil
}

func (r *PgTelegramNotificationRepository) FindOne(criteria types.Criteria) (types.TelegramNotification, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types.TelegramNotification{}, exception.
			New(err.Error()).
			Trace("r.transformer.Transform", "pg-telegram-notification-repository.go")
	}

	var (
		id            string
		session_id    string
		business_id   string
		booking_id    string
		scheduled_at  string
		chat_id       int
		business_name string
		created_at    string
	)

	if scanErr := r.connection.QueryRow(query).Scan(
		&id,
		&session_id,
		&business_id,
		&booking_id,
		&scheduled_at,
		&chat_id,
		&business_name,
		&created_at,
	); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return types.TelegramNotification{}, exception.
				New("Telegram notification not found").
				Trace("r.connection.QueryRow", "pg-telegram-notification-repository.go").
				WithValues([]string{query}).
				Domain()
		}

		return types.TelegramNotification{}, exception.
			New(scanErr.Error()).
			Trace("r.connection.QueryRow", "pg-telegram-notification-repository.go").
			WithValues([]string{query})
	}

	return types.TelegramNotification{
		Id:          id,
		SessionId:   session_id,
		BusinessId:  business_id,
		BookingId:   booking_id,
		ScheduledAt: scheduled_at,
		ChatId:      chat_id,
		From:        business_name,
		CreatedAt:   created_at,
	}, nil
}

func (r *PgTelegramNotificationRepository) Save(entity types.TelegramNotification) error {
	var query strings.Builder

	query.WriteString(`INSERT INTO telegram_notification `)
	query.WriteString(`(id, session_id, business_id, booking_id, scheduled_at, chat_id, business_name, created_at) `)
	query.WriteString(`VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`)

	_, err := r.connection.Exec(
		query.String(),
		entity.Id,
		entity.SessionId,
		entity.BusinessId,
		entity.BookingId,
		entity.ScheduledAt,
		entity.ChatId,
		entity.From,
		entity.CreatedAt,
	)

	if err != nil {
		return exception.
			New(err.Error()).
			Trace("r.connection.Exec", "pg-telegram-notification-repository.go").
			WithValues([]string{
				query.String(),
				entity.Id,
				entity.BusinessId,
				entity.BookingId,
				entity.ScheduledAt,
				entity.From,
				entity.CreatedAt,
			})
	}

	return nil
}

func (r *PgTelegramNotificationRepository) Update(_ types.TelegramNotification) error {
	return exception.
		New("Method not implemented").
		Trace("Update", "pg-telegram-notification-repository.go")
}
