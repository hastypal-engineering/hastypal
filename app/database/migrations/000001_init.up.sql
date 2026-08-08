/*
================================================================================
TABLES
================================================================================
*/

CREATE TABLE IF NOT EXISTS ha_business (
    hab_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hab_name VARCHAR(255) NOT NULL,
    hab_contact_phone VARCHAR(36) NOT NULL,
    hab_email VARCHAR(60) NOT NULL,
    hab_address VARCHAR(255) NOT NULL,
    hab_country VARCHAR(3) NOT NULL,
    hab_lang VARCHAR(3),
    hab_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    hab_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS ha_service_catalog (
    hasc_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hasc_name VARCHAR(255) NOT NULL,
    hasc_price INTEGER NOT NULL,
    hasc_currency VARCHAR(10) NOT NULL,
    hasc_duration VARCHAR(10),
    hasc_business_id BIGINT,
    hasc_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    hasc_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT fk_service_catalog_business FOREIGN KEY(hasc_business_id) REFERENCES ha_business(hab_id)
);

CREATE TABLE IF NOT EXISTS ha_employees (
    hae_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hae_name VARCHAR(255) NOT NULL,
    hae_business_id BIGINT NOT NULL,
    hae_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    hae_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT fk_employees_business FOREIGN KEY(hae_business_id) REFERENCES ha_business(hab_id)
)

CREATE TABLE IF NOT EXISTS ha_open_hours (
    haoh_business_id BIGINT,
    haoh_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    haoh_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT fk_service_catalog_business FOREIGN KEY(haoh_business_id) REFERENCES ha_business(hab_id)
)

CREATE TABLE IF NOT EXISTS ha_business_holidays (
    habh_business_id BIGINT,
    habh_date_add TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    habh_date_upd TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT fk_service_catalog_business FOREIGN KEY(habh_business_id) REFERENCES ha_business(hab_id)
)

CREATE TABLE IF NOT EXISTS booking_session (
    id VARCHAR(36) PRIMARY KEY, -- Unique identifier (UUID)
    business_id VARCHAR(36) NOT NULL,
    chat_id VARCHAR(36) NOT NULL,
    service_id VARCHAR(36) NOT NULL,
    date VARCHAR(60) NOT NULL,
    hour VARCHAR(5) NOT NULL,
    created_at VARCHAR(60) NOT NULL,
    updated_at VARCHAR(60) NOT NULL,
    ttl INTEGER NOT NULL,
    FOREIGN KEY (business_id) REFERENCES business(id)
);

CREATE TABLE IF NOT EXISTS booking (
    id VARCHAR(36) PRIMARY KEY, -- Unique identifier (UUID)
    session_id VARCHAR(8) NOT NULL,
    business_id VARCHAR(36) NOT NULL,
    service_id VARCHAR(36) NOT NULL,
    booking_date VARCHAR(60) NOT NULL,
    created_at VARCHAR(60) NOT NULL,
    FOREIGN KEY (session_id) REFERENCES booking_session(id),
    FOREIGN KEY (business_id) REFERENCES business(id)
);

CREATE TABLE IF NOT EXISTS telegram_notification (
    id VARCHAR(36) PRIMARY KEY, -- Unique identifier (UUID)
    session_id VARCHAR(8) NOT NULL,
    booking_id VARCHAR(36) NOT NULL,
    business_id VARCHAR(36) NOT NULL,
    scheduled_at VARCHAR(60) NOT NULL,
    chat_id INT NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    booking_date VARCHAR(60) NOT NULL,
    sent BOOLEAN NOT NULL,
    sent_at VARCHAR(60) NULL,
    created_at VARCHAR(60) NOT NULL,
    FOREIGN KEY (session_id) REFERENCES booking_session(id),
    FOREIGN KEY (booking_id) REFERENCES booking(id),
    FOREIGN KEY (business_id) REFERENCES business(id)
);

CREATE TABLE IF NOT EXISTS google_token (
    business_id VARCHAR(36) PRIMARY KEY, -- Unique identifier (UUID)
    access_token VARCHAR(255) NOT NULL,
    token_type VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    created_at VARCHAR(60) NOT NULL,
    updated_at VARCHAR(60) NOT NULL,
    FOREIGN KEY (business_id) REFERENCES business(id)
);