CREATE TABLE IF NOT EXISTS redirects (
    id SERIAL PRIMARY KEY,
    original TEXT NOT NULL UNIQUE,
    short VARCHAR(7) NOT NULL,
    user_ip VARCHAR(40),
    os VARCHAR(80),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
)