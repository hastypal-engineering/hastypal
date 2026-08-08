CREATE TABLE IF NOT EXISTS business (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(36) NOT NULL,
    email VARCHAR(60) NOT NULL,
    password VARCHAR(255) NOT NULL,
    channel_name VARCHAR(255),
    location VARCHAR(36),
    opening_hours JSON,
    holidays JSON,
    created_at VARCHAR(60) NOT NULL,
    updated_at VARCHAR(60) NOT NULL
);
CREATE TABLE IF NOT EXISTS service_catalog (
    id VARCHAR(36) PRIMARY key,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    currency VARCHAR(10) NOT NULL,
    duration VARCHAR(10),
    business_id VARCHAR(36),
    CONSTRAINT fk_business FOREIGN KEY(business_id) REFERENCES business(id)
);

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