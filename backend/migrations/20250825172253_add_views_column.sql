-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE pastes ADD COLUMN views INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE pastes DROP COLUMN views;
-- +goose StatementEnd
