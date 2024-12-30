ALTER TABLE business
    RENAME COLUMN country to location,
    DROP COLUMN street,
    DROP COLUMN post_code,
    DROP COLUMN city;