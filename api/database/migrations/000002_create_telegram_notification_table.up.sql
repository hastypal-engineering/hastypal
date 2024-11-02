CREATE TABLE IF NOT EXISTS telegram_notification (
    id VARCHAR(36) PRIMARY KEY, -- Unique identifier (UUID)
    session_id VARCHAR(8) NOT NULL,
    scheduled_at VARCHAR(60) NOT NULL,
    chat_id INT NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    FOREIGN KEY (session_id) REFERENCES booking_session(id),
    created_at VARCHAR(60) NOT NULL
);