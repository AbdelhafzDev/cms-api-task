DROP TRIGGER IF EXISTS trg_program_soft_delete ON programs;
DROP FUNCTION IF EXISTS notify_program_soft_delete();

-- Restore the ON DELETE CASCADE FK
ALTER TABLE search_index_jobs
    ADD CONSTRAINT search_index_jobs_program_id_fkey
    FOREIGN KEY (program_id) REFERENCES programs(id) ON DELETE CASCADE
    NOT VALID;
