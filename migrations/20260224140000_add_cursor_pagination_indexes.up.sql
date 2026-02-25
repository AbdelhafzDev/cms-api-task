-- Composite index for admin program list: keyset pagination on (created_at, id)
CREATE INDEX idx_programs_created_at_id ON programs (created_at DESC, id DESC);

-- Partial composite index for public discovery list: keyset pagination on (published_at, id)
CREATE INDEX idx_programs_published_at_id ON programs (published_at DESC, id DESC) WHERE status = 'published';
