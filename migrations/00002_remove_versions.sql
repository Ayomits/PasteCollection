
-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

ALTER TABLE pastes
DROP COLUMN version,
DROP COLUMN latest;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';

ALTER TABLE pastes
ADD COLUMN version VARCHAR(32),
ADD COLUMN latest BOOLEAN;
-- +goose StatementEnd
