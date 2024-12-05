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
    created_at VARCHAR(60) NOT NULL,
    FOREIGN KEY (session_id) REFERENCES booking_session(id),
    FOREIGN KEY (booking_id) REFERENCES booking(id),
    FOREIGN KEY (business_id) REFERENCES business(id)
);