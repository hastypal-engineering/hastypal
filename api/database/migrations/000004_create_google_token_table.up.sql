CREATE TABLE IF NOT EXISTS google_token (
    business_id VARCHAR(36) PRIMARY KEY, -- Unique identifier (UUID)
    access_token VARCHAR(255) NOT NULL,
    token_type VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    created_at VARCHAR(60) NOT NULL,
    updated_at VARCHAR(60) NOT NULL
    -- FOREIGN KEY (business_id) REFERENCES business(id)
);