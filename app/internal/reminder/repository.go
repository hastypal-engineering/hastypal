package reminder

import (
	"context"
	"database/sql"
)

type ReminderRepository interface {
	Save(ctx context.Context, reminder *Reminder) error
}

type PgReminderRepository struct {
	connection *sql.DB
}

func NewPgReminderRepository(connection *sql.DB) *PgReminderRepository {
	return &PgReminderRepository{
		connection: connection,
	}
}

func (r *PgReminderRepository) Save(ctx context.Context, reminder *Reminder) error {
	return nil
}
