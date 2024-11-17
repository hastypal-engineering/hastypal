package repository

import (
	"database/sql"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
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

func (r *PgTelegramNotificationRepository) Find(criteria types2.Criteria) ([]types2.TelegramNotification, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return nil, types2.ApiError{
			Msg:      err.Error(),
			Function: "Find -> r.transformer.Transform()",
			File:     "pg-telegram-notification-repository.go",
		}
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, types2.ApiError{
			Msg:      queryErr.Error(),
			Function: "Find -> r.connection.Query()",
			File:     "pg-telegram-notification-repository.go",
			Values:   []string{query},
		}
	}

	defer rows.Close()

	var (
		id            string
		session_id    string
		scheduled_at  string
		chat_id       int
		business_name string
		created_at    string
	)

	var results []types2.TelegramNotification

	for rows.Next() {
		if scanErr := rows.Scan(
			&id,
			&session_id,
			&scheduled_at,
			&chat_id,
			&business_name,
			&created_at,
		); scanErr != nil {
			return nil, types2.ApiError{
				Msg:      scanErr.Error(),
				Function: "Find -> rows.Scan()",
				File:     "pg-telegram-notification-repository.go",
			}
		}

		results = append(results, types2.TelegramNotification{
			Id:          id,
			SessionId:   session_id,
			ScheduledAt: scheduled_at,
			ChatId:      chat_id,
			From:        business_name,
			CreatedAt:   created_at,
		})
	}

	return results, nil
}

func (r *PgTelegramNotificationRepository) FindOne(criteria types2.Criteria) (types2.TelegramNotification, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types2.TelegramNotification{}, types2.ApiError{
			Msg:      err.Error(),
			Function: "FindOne -> r.transformer.Transform()",
			File:     "pg-telegram-notification-repository.go",
		}
	}

	var (
		id            string
		session_id    string
		scheduled_at  string
		chat_id       int
		business_name string
		created_at    string
	)

	if scanErr := r.connection.QueryRow(query).Scan(
		&id,
		&session_id,
		&scheduled_at,
		&chat_id,
		&business_name,
		&created_at,
	); scanErr != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return types2.TelegramNotification{}, types2.ApiError{
				Msg:      "Entity Business not found",
				Function: "FindOne -> rows.Scan()",
				File:     "pg-telegram-notification-repository.go",
				Values:   []string{query},
				Domain:   true,
			}
		}

		return types2.TelegramNotification{}, types2.ApiError{
			Msg:      scanErr.Error(),
			Function: "FindOne -> rows.Scan()",
			File:     "pg-telegram-notification-repository.go",
			Values:   []string{query},
		}
	}

	return types2.TelegramNotification{
		Id:          id,
		SessionId:   session_id,
		ScheduledAt: scheduled_at,
		ChatId:      chat_id,
		From:        business_name,
		CreatedAt:   created_at,
	}, nil
}

func (r *PgTelegramNotificationRepository) Save(entity types2.TelegramNotification) error {
	var query strings.Builder

	query.WriteString(`INSERT INTO telegram_notification `)
	query.WriteString(`(id, session_id, scheduled_at, chat_id, business_name, created_at) `)
	query.WriteString(`VALUES ($1, $2, $3, $4, $5, $6);`)

	_, err := r.connection.Exec(
		query.String(),
		entity.Id,
		entity.SessionId,
		entity.ScheduledAt,
		entity.ChatId,
		entity.From,
		entity.CreatedAt,
	)

	if err != nil {
		return types2.ApiError{
			Msg:      err.Error(),
			Function: "Save -> r.connection.Exec()",
			File:     "pg-telegram-notification-repository.go",
			Values: []string{
				query.String(),
				entity.Id,
				entity.ScheduledAt,
				entity.From,
				entity.CreatedAt,
			},
		}
	}

	return nil
}

func (r *PgTelegramNotificationRepository) Update(_ types2.TelegramNotification) error {
	return types2.ApiError{
		Msg:      "Method not implemented yet",
		Function: "Update",
		File:     "pg-telegram-notification-repository.go",
	}
}