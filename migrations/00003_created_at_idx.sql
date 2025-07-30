-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DROP INDEX idx_pastes_id;
CREATE INDEX idx_pastes_created_at_id ON pastes(created_at DESC, id DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_pastes_created_at_id;
CREATE UNIQUE INDEX IF NOT EXISTS idx_pastes_id ON pastes (id);
-- +goose StatementEnd
