-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pastes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(32) UNIQUE NOT NULL,
    tags VARCHAR(32)[],
    paste TEXT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_pastes_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);


CREATE UNIQUE INDEX IF NOT EXISTS idx_pastes_id ON pastes (id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_pastes_title ON pastes(title);
CREATE INDEX IF NOT EXISTS idx_pastes_text ON pastes(paste);
CREATE INDEX IF NOT EXISTS idx_pastes_user ON pastes(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_pastes_id;
DROP INDEX IF EXISTS idx_pastes_text;
DROP INDEX IF EXISTS idx_pastes_title;
DROP INDEX IF EXISTS idx_pastes_user;
DROP TABLE IF EXISTS pastes;
-- +goose StatementEnd
