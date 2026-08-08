CREATE TABLE IF NOT EXISTS urls (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    long_url TEXT NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_urls_expired_at ON urls(expired_at);