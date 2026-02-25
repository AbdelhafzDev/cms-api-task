-- Drop the ON DELETE CASCADE FK so delete index jobs survive program deletion.
-- Queue/outbox tables should not have cascading FK constraints.
ALTER TABLE search_index_jobs
    DROP CONSTRAINT search_index_jobs_program_id_fkey;

-- Trigger function: enqueue a delete job when a program is soft-deleted.
-- Fires BEFORE UPDATE when deleted_at transitions from NULL to non-NULL.
CREATE OR REPLACE FUNCTION notify_program_soft_delete() RETURNS TRIGGER AS $$
BEGIN
    IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
        -- Remove completed/dead jobs for this program (no longer needed)
        DELETE FROM search_index_jobs
        WHERE program_id = OLD.id
          AND status IN ('completed', 'dead');

        -- Enqueue the delete job (dedupe with active partial unique index)
        INSERT INTO search_index_jobs (program_id, action, status, scheduled_at)
        VALUES (OLD.id, 'delete', 'pending', NOW())
        ON CONFLICT (program_id, action) WHERE status IN ('pending', 'processing', 'failed')
        DO UPDATE SET scheduled_at = NOW(), updated_at = NOW();

        -- Cancel any pending upsert job for this program
        DELETE FROM search_index_jobs
        WHERE program_id = OLD.id
          AND action = 'upsert'
          AND status IN ('pending', 'failed');
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_program_soft_delete
    BEFORE UPDATE ON programs
    FOR EACH ROW
    EXECUTE FUNCTION notify_program_soft_delete();
