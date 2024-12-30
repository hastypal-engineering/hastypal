ALTER TABLE business
    RENAME COLUMN location TO country;
ALTER TABLE business
ADD street VARCHAR(36);
ALTER TABLE business
ADD post_code VARCHAR(36);
ALTER TABLE business
ADD city VARCHAR(36);