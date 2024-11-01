package repository

import (
	"database/sql"
	"errors"
	"github.com/adriein/hastypal/internal/hastypal/helper"
	"github.com/adriein/hastypal/internal/hastypal/types"
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
		return nil, types.ApiError{
			Msg:      err.Error(),
			Function: "Find -> r.transformer.Transform()",
			File:     "pg-telegram-notification-repository.go",
		}
	}

	rows, queryErr := r.connection.Query(query)

	if queryErr != nil {
		return nil, types.ApiError{
			Msg:      queryErr.Error(),
			Function: "Find -> r.connection.Query()",
			File:     "pg-telegram-notification-repository.go",
			Values:   []string{query},
		}
	}

	defer rows.Close()

	var (
		id           string
		scheduled_at string
		chat_id      int
		from         string
		created_at   string
	)

	var results []types.TelegramNotification

	for rows.Next() {
		if scanErr := rows.Scan(
			&id,
			&scheduled_at,
			&chat_id,
			&from,
			&created_at,
		); scanErr != nil {
			return nil, types.ApiError{
				Msg:      scanErr.Error(),
				Function: "Find -> rows.Scan()",
				File:     "pg-telegram-notification-repository.go",
			}
		}

		results = append(results, types.TelegramNotification{
			Id:          id,
			ScheduledAt: scheduled_at,
			ChatId:      chat_id,
			From:        from,
			CreatedAt:   created_at,
		})
	}

	return results, nil
}

func (r *PgTelegramNotificationRepository) FindOne(criteria types.Criteria) (types.TelegramNotification, error) {
	query, err := r.transformer.Transform(criteria)

	if err != nil {
		return types.TelegramNotification{}, types.ApiError{
			Msg:      err.Error(),
			Function: "FindOne -> r.transformer.Transform()",
			File:     "pg-telegram-notification-repository.go",
		}
	}

	var (
		id           string
		scheduled_at string
		chat_id      int
		from         string
		created_at   string
	)

	if scanErr := r.connection.QueryRow(query).Scan(
		&id,
		&scheduled_at,
		&chat_id,
		&from,
		&created_at,
	); scanErr != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return types.TelegramNotification{}, types.ApiError{
				Msg:      "Entity Business not found",
				Function: "FindOne -> rows.Scan()",
				File:     "pg-telegram-notification-repository.go",
				Values:   []string{query},
				Domain:   true,
			}
		}

		return types.TelegramNotification{}, types.ApiError{
			Msg:      scanErr.Error(),
			Function: "FindOne -> rows.Scan()",
			File:     "pg-telegram-notification-repository.go",
			Values:   []string{query},
		}
	}

	return types.TelegramNotification{
		Id:          id,
		ScheduledAt: scheduled_at,
		ChatId:      chat_id,
		From:        from,
		CreatedAt:   created_at,
	}, nil
}

func (r *PgTelegramNotificationRepository) Save(entity types.TelegramNotification) error {
	var query = `INSERT INTO notification (id, scheduled_at, chat_id, from, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.connection.Exec(
		query,
		entity.Id,
		entity.ScheduledAt,
		entity.ChatId,
		entity.From,
		entity.CreatedAt,
	)

	if err != nil {
		return types.ApiError{
			Msg:      err.Error(),
			Function: "Save -> r.connection.Exec()",
			File:     "pg-telegram-notification-repository.go",
			Values: []string{
				query,
				entity.Id,
				entity.ScheduledAt,
				entity.From,
				entity.CreatedAt,
			},
		}
	}

	return nil
}

func (r *PgTelegramNotificationRepository) Update(_ types.TelegramNotification) error {
	return types.ApiError{
		Msg:      "Method not implemented yet",
		Function: "Update",
		File:     "pg-telegram-notification-repository.go",
	}
}
