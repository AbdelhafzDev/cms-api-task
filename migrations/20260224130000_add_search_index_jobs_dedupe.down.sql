-- Restore the original trigger function (no-op ON CONFLICT DO NOTHING)
CREATE OR REPLACE FUNCTION notify_program_index() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO search_index_jobs (program_id, action, status, scheduled_at)
    VALUES (NEW.id, 'upsert', 'pending', NOW())
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop the partial unique index
DROP INDEX IF EXISTS idx_search_index_jobs_active;
