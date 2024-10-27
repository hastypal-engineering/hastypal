CREATE TABLE IF NOT EXISTS business (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(36) NOT NULL,
    email VARCHAR(60) NOT NULL,
    password VARCHAR(255) NOT NULL,
    channel_name VARCHAR(255),
    location VARCHAR(36),
    opening_hours JSON,
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