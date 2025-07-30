
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
ADD COLUMN version VARCHAR(32) DEFAULT "1.0.0",
ADD COLUMN latest BOOLEAN DEFAULT true;
-- +goose StatementEnd
