CREATE TABLE search_index_jobs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id    UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    action        VARCHAR(10) NOT NULL DEFAULT 'upsert',
    status        VARCHAR(15) NOT NULL DEFAULT 'pending',
    attempts      INT NOT NULL DEFAULT 0,
    max_attempts  INT NOT NULL DEFAULT 5,
    last_error    TEXT,
    scheduled_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_action CHECK (action IN ('upsert', 'delete')),
    CONSTRAINT chk_status CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'dead'))
);

CREATE INDEX idx_search_index_jobs_poll ON search_index_jobs(status, scheduled_at);
CREATE INDEX idx_search_index_jobs_program ON search_index_jobs(program_id);

-- Trigger function: auto-create an index job when a program is inserted or updated
CREATE OR REPLACE FUNCTION notify_program_index() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO search_index_jobs (program_id, action, status, scheduled_at)
    VALUES (NEW.id, 'upsert', 'pending', NOW())
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_program_index
    AFTER INSERT OR UPDATE ON programs
    FOR EACH ROW
    EXECUTE FUNCTION notify_program_index();

-- Create index jobs for existing programs
INSERT INTO search_index_jobs (program_id, action, status, scheduled_at)
SELECT id, 'upsert', 'pending', NOW()
FROM programs;
