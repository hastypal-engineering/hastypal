/*
================================================================================
TABLES
================================================================================
*/

CREATE TABLE IF NOT EXISTS ha_business (
    hab_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hab_public_id VARCHAR(255) NOT NULL,
    hab_name VARCHAR(255) NOT NULL,
    hab_contact_phone VARCHAR(36) NOT NULL,
    hab_email VARCHAR(60) NOT NULL,
    hab_address VARCHAR(255) NOT NULL,
    hab_country VARCHAR(3) NOT NULL,
    hab_lang VARCHAR(3) NOT NULL,
    hab_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    hab_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS ha_service_catalog (
    hasc_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hasc_name VARCHAR(255) NOT NULL,
    hasc_description VARCHAR(255) NOT NULL,
    hasc_price NUMERIC(12, 2) NOT NULL,
    hasc_currency VARCHAR(10) NOT NULL,
    hasc_duration VARCHAR(10) NOT NULL,
    hasc_business_id BIGINT NOT NULL,
    hasc_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    hasc_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT fk_service_catalog_business FOREIGN KEY(hasc_business_id) 
      REFERENCES ha_business(hab_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ha_booking_session (
    habs_id VARCHAR(36) PRIMARY KEY,
    habs_business_id BIGINT NOT NULL,
    habs_service_id BIGINT NULL,
    habs_date DATE NULL,
    habs_time TIME NULL,
    habs_ttl INTEGER NOT NULL,
    habs_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    habs_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT fk_booking_session_business FOREIGN KEY(habs_business_id) 
      REFERENCES ha_business(hab_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ha_business_holiday (
    habh_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    habh_business_id BIGINT NOT NULL,
    habh_name VARCHAR NOT NULL,
    habh_start_date DATE NOT NULL,
    habh_end_date DATE NOT NULL,
    habh_is_recurring BOOLEAN NOT NULL,
    habh_is_closed BOOLEAN NOT NULL,
    habh_type VARCHAR(36) NOT NULL,
    habh_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    habh_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT fk_holiday_business FOREIGN KEY(habh_business_id) 
      REFERENCES ha_business(hab_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ha_business_operating_day (
    habod_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    habod_business_id BIGINT NOT NULL,
    habod_day_of_week SMALLINT NOT NULL,
    habod_is_closed BOOLEAN NOT NULL,
    CONSTRAINT fk_operating_day_business FOREIGN KEY(habod_business_id) 
      REFERENCES ha_business(hab_id) ON DELETE CASCADE,
    CONSTRAINT uq_business_operating_day UNIQUE (habod_business_id, habod_day_of_week)
);

CREATE TABLE IF NOT EXISTS ha_business_schedule_override (
    habso_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    habso_business_id BIGINT NOT NULL,
    habso_date DATE NOT NULL,
    habso_is_closed BOOLEAN NOT NULL,
    habso_reason VARCHAR NOT NULL,
    CONSTRAINT fk_schedule_override_business FOREIGN KEY(habso_business_id) 
      REFERENCES ha_business(hab_id) ON DELETE CASCADE,
    CONSTRAINT uq_business_override_date UNIQUE (habso_business_id, habso_date)
);

CREATE TABLE IF NOT EXISTS ha_business_time_slot (
    habts_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    habts_operating_day_id BIGINT NULL,
    habts_override_id BIGINT NULL,
    habts_open_time TIME NOT NULL,
    habts_close_time TIME NOT NULL,
    CONSTRAINT fk_time_slot_operating_day FOREIGN KEY(habts_operating_day_id) 
      REFERENCES ha_business_operating_day(habod_id) ON DELETE CASCADE,
    CONSTRAINT fk_time_slot_override FOREIGN KEY(habts_override_id)
      REFERENCES ha_business_schedule_override(habso_id) ON DELETE CASCADE,
    CONSTRAINT ck_time_slot_parent CHECK (
      (habts_operating_day_id IS NOT NULL AND habts_override_id IS NULL) OR
      (habts_operating_day_id IS NULL AND habts_override_id IS NOT NULL)
    )
);
-- CREATE TABLE IF NOT EXISTS ha_employees (
--     hae_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     hae_name VARCHAR(255) NOT NULL,
--     hae_business_id BIGINT NOT NULL,
--     hae_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
--     hae_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
--     CONSTRAINT fk_employees_business FOREIGN KEY(hae_business_id) REFERENCES ha_business(hab_id)
-- )

--
-- CREATE TABLE IF NOT EXISTS booking (
--     id VARCHAR(36) PRIMARY KEY,
--     session_id VARCHAR(8) NOT NULL,
--     business_id VARCHAR(36) NOT NULL,
--     service_id VARCHAR(36) NOT NULL,
--     booking_date VARCHAR(60) NOT NULL,
--     created_at VARCHAR(60) NOT NULL,
--     FOREIGN KEY (session_id) REFERENCES booking_session(id),
--     FOREIGN KEY (business_id) REFERENCES business(id)
-- );
--
-- CREATE TABLE IF NOT EXISTS telegram_notification (
--     id VARCHAR(36) PRIMARY KEY,
--     session_id VARCHAR(8) NOT NULL,
--     booking_id VARCHAR(36) NOT NULL,
--     business_id VARCHAR(36) NOT NULL,
--     scheduled_at VARCHAR(60) NOT NULL,
--     chat_id INT NOT NULL,
--     business_name VARCHAR(255) NOT NULL,
--     service_name VARCHAR(255) NOT NULL,
--     booking_date VARCHAR(60) NOT NULL,
--     sent BOOLEAN NOT NULL,
--     sent_at VARCHAR(60) NULL,
--     created_at VARCHAR(60) NOT NULL,
--     FOREIGN KEY (session_id) REFERENCES booking_session(id),
--     FOREIGN KEY (booking_id) REFERENCES booking(id),
--     FOREIGN KEY (business_id) REFERENCES business(id)
-- );
--
-- CREATE TABLE IF NOT EXISTS google_token (
--     business_id VARCHAR(36) PRIMARY KEY, -- Unique identifier (UUID)
--     access_token VARCHAR(255) NOT NULL,
--     token_type VARCHAR(255) NOT NULL,
--     refresh_token VARCHAR(255) NOT NULL,
--     created_at VARCHAR(60) NOT NULL,
--     updated_at VARCHAR(60) NOT NULL,
--     FOREIGN KEY (business_id) REFERENCES business(id)
-- );
