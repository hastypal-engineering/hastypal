package business

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rotisserie/eris"
)

var BusinessNotFound = eris.New("Business not found")

type BusinessRepository interface {
	GetByID(ctx context.Context, ID int) (*Business, error)
}

type PgBusinessRepository struct {
	connection *sql.DB
}

func NewPgBusinessRepository(connection *sql.DB) *PgBusinessRepository {
	return &PgBusinessRepository{
		connection: connection,
	}
}

func (r *PgBusinessRepository) GetByID(ctx context.Context, ID int) (*Business, error) {
	query := `
		SELECT
			a.ta_id,
			a.ta_auction_id,
			a.ta_tibia_auction_link,
			a.ta_img,
			a.ta_char_name,
			a.ta_char_level,
			v.*,
			g.*,
			w.*,
			ts.*,
			a.ta_world_transfer,
			a.ta_boss_points,
			a.ta_charm_expansion,
			a.ta_charm_points,
			a.ta_task_expansion,
			a.ta_current_bid,
			a.ta_current_bid_fiat,
			a.ta_current_bid_currency,
			a.ta_auction_stage,
			a.ta_auction_start,
			a.ta_auction_end,
			tar.tar_status,
			tar.tar_date_add,
			tar.tar_date_upd
		FROM
			tc_auction a
		INNER JOIN
			tc_vocation v ON a.ta_char_vocation = v.tv_id
		INNER JOIN
			tc_gender g ON a.ta_char_gender = g.tg_id
		INNER JOIN
			tc_world w ON a.ta_char_world = w.tw_id
		INNER JOIN
			tc_auction_recording tar ON a.ta_id = tar.tar_recordable_id
		INNER JOIN
			tc_skills ts ON a.ta_auction_id = ts.ts_auction_id
		WHERE
			a.ta_auction_id = $1;
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var (
		business Business
	)

	err := r.connection.QueryRowContext(ctxTimeout, query, ID).Scan(
		&business.ID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, BusinessNotFound
		}

		return nil, eris.Wrap(err, "Failed to query business by ID")
	}

	return &business, nil
}
