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