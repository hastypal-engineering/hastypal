package business

import (
	"context"
	"database/sql"
	"time"

	"github.com/adriein/hastypal/database"
	"github.com/rotisserie/eris"
)

var BusinessNotFound = eris.New("Business not found")

const timeOnlyFormat = "15:04"

type BusinessRepository interface {
	Create(ctx context.Context, business *Business) (int, error)
	CreateService(ctx context.Context, service *ServiceCatalog) (int, error)
	CreateSchedule(ctx context.Context, businessID int, schedule *BusinessSchedule) error
	GetSchedule(ctx context.Context, businessID int) (*BusinessSchedule, error)
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

func (r *PgBusinessRepository) CreateSchedule(ctx context.Context, businessID int, schedule *BusinessSchedule) error {
	operatingDayQuery := `
		INSERT INTO ha_business_operating_day (
			habod_business_id,
			habod_day_of_week,
			habod_is_closed
		)
		VALUES ($1, $2, $3)
		RETURNING habod_id;
	`

	timeSlotQuery := `
		INSERT INTO ha_business_time_slot (
			habts_operating_day_id,
			habts_open_time,
			habts_close_time
		)
		VALUES ($1, $2, $3);
	`

	holidayQuery := `
		INSERT INTO ha_business_holiday (
			habh_business_id,
			habh_name,
			habh_start_date,
			habh_end_date,
			habh_is_recurring,
			habh_is_closed,
			habh_type,
			habh_date_add,
			habh_date_upd
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`

	overrideQuery := `
		INSERT INTO ha_business_schedule_override (
			habso_business_id,
			habso_date,
			habso_is_closed,
			habso_reason
		)
		VALUES ($1, $2, $3, $4)
		RETURNING habso_id;
	`

	overrideTimeSlotQuery := `
		INSERT INTO ha_business_time_slot (
			habts_override_id,
			habts_open_time,
			habts_close_time
		)
		VALUES ($1, $2, $3);
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	now := time.Now()

	tx, err := r.connection.BeginTx(ctxTimeout, nil)
	if err != nil {
		return eris.Wrap(err, "Failed to start schedule transaction")
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = eris.Wrapf(err, "Failed to rollback schedule transaction: %s", rollbackErr)
			}
		}
	}()

	for _, day := range schedule.WeeklySchedule {
		var operatingDayID int

		err = tx.QueryRowContext(ctxTimeout, operatingDayQuery,
			businessID,
			int(day.DayOfWeek),
			day.IsClosed,
		).Scan(&operatingDayID)
		if err != nil {
			return eris.Wrapf(err, "Failed to create operating day %s", day.DayOfWeek)
		}

		for _, slot := range day.TimeSlots {
			openTime, closeTime, parseErr := parseTimeSlot(slot)
			if parseErr != nil {
				return parseErr
			}

			if _, err = tx.ExecContext(ctxTimeout, timeSlotQuery, operatingDayID, openTime, closeTime); err != nil {
				return eris.Wrapf(err, "Failed to create time slot %s-%s for operating day %d", slot.OpenTime, slot.CloseTime, operatingDayID)
			}
		}
	}

	for _, holiday := range schedule.Holidays {
		if _, err = tx.ExecContext(ctxTimeout, holidayQuery,
			businessID,
			holiday.Name,
			holiday.StartDate,
			holiday.EndDate,
			holiday.IsRecurring,
			holiday.IsClosed,
			holiday.Type,
			now,
			now,
		); err != nil {
			return eris.Wrapf(err, "Failed to create holiday %s", holiday.Name)
		}
	}

	for _, override := range schedule.Overrides {
		var overrideID int

		err = tx.QueryRowContext(ctxTimeout, overrideQuery,
			businessID,
			override.Date,
			override.IsClosed,
			override.Reason,
		).Scan(&overrideID)
		if err != nil {
			return eris.Wrapf(err, "Failed to create schedule override on %s", override.Date)
		}

		for _, slot := range override.TimeSlots {
			openTime, closeTime, parseErr := parseTimeSlot(slot)
			if parseErr != nil {
				return parseErr
			}

			if _, err = tx.ExecContext(ctxTimeout, overrideTimeSlotQuery, overrideID, openTime, closeTime); err != nil {
				return eris.Wrapf(err, "Failed to create time slot %s-%s for schedule override %d", slot.OpenTime, slot.CloseTime, overrideID)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return eris.Wrap(err, "Failed to commit schedule transaction")
	}

	return nil
}

func (r *PgBusinessRepository) GetSchedule(ctx context.Context, businessID int) (*BusinessSchedule, error) {
	weeklyScheduleQuery := `
		SELECT
			d.habod_day_of_week,
			d.habod_is_closed,
			s.habts_open_time,
			s.habts_close_time
		FROM
			ha_business_operating_day d
		LEFT JOIN
			ha_business_time_slot s ON s.habts_operating_day_id = d.habod_id
		WHERE
			d.habod_business_id = $1
		ORDER BY
			d.habod_day_of_week,
			s.habts_open_time;
	`

	overrideQuery := `
		SELECT
			o.habso_id,
			o.habso_date,
			o.habso_is_closed,
			o.habso_reason,
			s.habts_open_time,
			s.habts_close_time
		FROM
			ha_business_schedule_override o
		LEFT JOIN
			ha_business_time_slot s ON s.habts_override_id = o.habso_id
		WHERE
			o.habso_business_id = $1
		ORDER BY
			o.habso_date,
			s.habts_open_time;
	`

	holidayQuery := `
		SELECT
			h.habh_id,
			h.habh_name,
			h.habh_start_date,
			h.habh_end_date,
			h.habh_is_recurring,
			h.habh_is_closed,
			h.habh_type
		FROM
			ha_business_holiday h
		WHERE
			h.habh_business_id = $1
		ORDER BY
			h.habh_start_date;
	`

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	schedule := &BusinessSchedule{}

	weeklyRows, err := r.connection.QueryContext(ctxTimeout, weeklyScheduleQuery, businessID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to query weekly schedule")
	}

	defer database.CloseRowsSafely(weeklyRows, &err)

	daysByWeekday := map[time.Weekday]*OperatingDay{}

	for weeklyRows.Next() {
		var (
			dayOfWeek int
			isClosed  bool
			openTime  sql.NullTime
			closeTime sql.NullTime
		)

		if err = weeklyRows.Scan(&dayOfWeek, &isClosed, &openTime, &closeTime); err != nil {
			return nil, eris.Wrap(err, "Failed to scan weekly schedule")
		}

		day, found := daysByWeekday[time.Weekday(dayOfWeek)]
		if !found {
			day = &OperatingDay{
				DayOfWeek: time.Weekday(dayOfWeek),
				IsClosed:  isClosed,
			}
			daysByWeekday[time.Weekday(dayOfWeek)] = day
			schedule.WeeklySchedule = append(schedule.WeeklySchedule, day)
		}

		if openTime.Valid && closeTime.Valid {
			day.TimeSlots = append(day.TimeSlots, TimeSlot{
				OpenTime:  openTime.Time.Format(timeOnlyFormat),
				CloseTime: closeTime.Time.Format(timeOnlyFormat),
			})
		}
	}

	if err = weeklyRows.Err(); err != nil {
		return nil, eris.Wrap(err, "Failed to iterate weekly schedule rows")
	}

	overrideRows, err := r.connection.QueryContext(ctxTimeout, overrideQuery, businessID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to query schedule overrides")
	}

	defer database.CloseRowsSafely(overrideRows, &err)

	overridesByID := map[int]*ScheduleOverride{}

	for overrideRows.Next() {
		var (
			overrideID int
			date       time.Time
			isClosed   bool
			reason     string
			openTime   sql.NullTime
			closeTime  sql.NullTime
		)

		if err = overrideRows.Scan(&overrideID, &date, &isClosed, &reason, &openTime, &closeTime); err != nil {
			return nil, eris.Wrap(err, "Failed to scan schedule overrides")
		}

		override, found := overridesByID[overrideID]
		if !found {
			override = &ScheduleOverride{
				Date:     date,
				IsClosed: isClosed,
				Reason:   reason,
			}
			overridesByID[overrideID] = override
			schedule.Overrides = append(schedule.Overrides, override)
		}

		if openTime.Valid && closeTime.Valid {
			override.TimeSlots = append(override.TimeSlots, TimeSlot{
				OpenTime:  openTime.Time.Format(timeOnlyFormat),
				CloseTime: closeTime.Time.Format(timeOnlyFormat),
			})
		}
	}

	if err = overrideRows.Err(); err != nil {
		return nil, eris.Wrap(err, "Failed to iterate schedule override rows")
	}

	holidayRows, err := r.connection.QueryContext(ctxTimeout, holidayQuery, businessID)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to query holidays")
	}

	defer database.CloseRowsSafely(holidayRows, &err)

	for holidayRows.Next() {
		var holiday Holiday

		if err = holidayRows.Scan(
			&holiday.ID,
			&holiday.Name,
			&holiday.StartDate,
			&holiday.EndDate,
			&holiday.IsRecurring,
			&holiday.IsClosed,
			&holiday.Type,
		); err != nil {
			return nil, eris.Wrap(err, "Failed to scan holidays")
		}

		schedule.Holidays = append(schedule.Holidays, &holiday)
	}

	if err = holidayRows.Err(); err != nil {
		return nil, eris.Wrap(err, "Failed to iterate holiday rows")
	}

	return schedule, nil
}

func parseTimeSlot(slot TimeSlot) (time.Time, time.Time, error) {
	openTime, err := time.Parse(timeOnlyFormat, slot.OpenTime)
	if err != nil {
		return time.Time{}, time.Time{}, eris.Wrapf(err, "Invalid time slot open time %s", slot.OpenTime)
	}

	closeTime, err := time.Parse(timeOnlyFormat, slot.CloseTime)
	if err != nil {
		return time.Time{}, time.Time{}, eris.Wrapf(err, "Invalid time slot close time %s", slot.CloseTime)
	}

	return openTime, closeTime, nil
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
