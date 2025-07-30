-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pastes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(32) UNIQUE NOT NULL,
    tags VARCHAR(32)[],
    paste TEXT NOT NULL,
    version VARCHAR(32) DEFAULT '1.0.0',
    latest BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_pastes_id ON pastes (id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_pastes_title ON pastes(title);
CREATE INDEX IF NOT EXISTS idx_pastes_text ON pastes(paste);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_pastes_id;
DROP INDEX IF EXISTS idx_pastes_text;
DROP INDEX IF EXISTS idx_pastes_title;
DROP TABLE IF EXISTS pastes;
-- +goose StatementEnd
