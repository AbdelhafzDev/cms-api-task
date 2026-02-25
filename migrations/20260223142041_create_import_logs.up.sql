CREATE TABLE import_logs (
    id               UUID PRIMARY KEY,
    source_id        BIGINT NOT NULL REFERENCES import_sources(id) ON DELETE CASCADE,
    triggered_by     UUID REFERENCES users(id) ON DELETE SET NULL,
    status           VARCHAR(15) NOT NULL DEFAULT 'pending',
    records_imported INTEGER NOT NULL DEFAULT 0,
    error_message    TEXT NOT NULL DEFAULT '',
    started_at       TIMESTAMPTZ,
    finished_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_import_logs_source_id ON import_logs(source_id);
