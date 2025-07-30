-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

CREATE TABLE
    IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(32) NOT NULL UNIQUE,
        social_id VARCHAR(255) NOT NULL UNIQUE
    );

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

CREATE INDEX IF NOT EXISTS idx_users_social_id ON users (social_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';

DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS idx_users_username;

DROP INDEX IF EXISTS idx_users_social_id;

-- +goose StatementEnd
