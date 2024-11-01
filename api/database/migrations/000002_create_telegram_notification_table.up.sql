CREATE TABLE IF NOT EXISTS telegram_notification (
    id VARCHAR(36) PRIMARY KEY, -- Unique identifier (UUID)
    scheduled_at VARCHAR(60) NOT NULL,
    chat_id INT NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    created_at VARCHAR(60) NOT NULL
);