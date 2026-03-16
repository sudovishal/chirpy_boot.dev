-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOL NOT NULL DEFAULT false;
-- ADD COLUMN hashed_password TEXT NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;
