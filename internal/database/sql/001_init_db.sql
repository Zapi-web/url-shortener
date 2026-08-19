-- +goose Up
CREATE TABLE IF NOT EXISTS urls (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    long_url TEXT NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id);
CREATE INDEX IF NOT EXISTS idx_urls_expired_at ON urls(expired_at);

-- +goose Down
DROP TABLE IF EXISTS urls;