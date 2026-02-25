CREATE TABLE import_sources (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    source_type VARCHAR(30) NOT NULL,
    base_url    VARCHAR(2048) NOT NULL DEFAULT '',
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed import sources
INSERT INTO import_sources (name, source_type, base_url) VALUES
    ('YouTube - Thmanyah', 'youtube', 'https://www.youtube.com/@thmanyahPodcasts');
