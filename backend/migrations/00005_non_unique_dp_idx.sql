-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP INDEX IF EXISTS idx_users_display_username;
CREATE INDEX IF NOT EXISTS idx_users_display_username on users(display_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_users_display_username;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_display_username on users(display_name);
-- +goose StatementEnd
