ALTER TABLE business
ADD street VARCHAR(36),
    ADD post_code VARCHAR(36),
    ADD city VARCHAR(36),
    ADD country VARCHAR(36),
    DROP COLUMN location;