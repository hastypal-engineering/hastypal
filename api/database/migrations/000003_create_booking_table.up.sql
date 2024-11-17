CREATE TABLE IF NOT EXISTS booking (
    id VARCHAR(36) PRIMARY KEY, -- Unique identifier (UUID)
    session_id VARCHAR(8) NOT NULL,
    business_id VARCHAR(36) NOT NULL,
    service_id VARCHAR(36) NOT NULL,
    booking_date VARCHAR(60) NOT NULL,
    created_at VARCHAR(60) NOT NULL,
    FOREIGN KEY (session_id) REFERENCES booking_session(id)
);