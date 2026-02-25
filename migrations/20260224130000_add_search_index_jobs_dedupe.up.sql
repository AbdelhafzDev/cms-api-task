-- Partial unique index: at most one active job per program/action
CREATE UNIQUE INDEX idx_search_index_jobs_active
    ON search_index_jobs (program_id, action)
    WHERE status IN ('pending', 'processing', 'failed');

-- Replace trigger function: ON CONFLICT resets scheduled_at so the job gets re-polled
CREATE OR REPLACE FUNCTION notify_program_index() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO search_index_jobs (program_id, action, status, scheduled_at)
    VALUES (NEW.id, 'upsert', 'pending', NOW())
    ON CONFLICT (program_id, action) WHERE status IN ('pending', 'processing', 'failed')
    DO UPDATE SET scheduled_at = NOW(), updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
