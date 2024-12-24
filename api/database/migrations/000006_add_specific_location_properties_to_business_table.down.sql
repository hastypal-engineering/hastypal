ALTER TABLE business
ADD location VARCHAR(36),
    DROP COLUMN street,
    DROP COLUMN post_code,
    DROP COLUMN city,
    DROP COLUMN country;