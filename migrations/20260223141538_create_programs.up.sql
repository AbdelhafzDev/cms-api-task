CREATE TABLE programs (
    id               UUID PRIMARY KEY,
    title            VARCHAR(255) NOT NULL,
    description      TEXT NOT NULL DEFAULT '',
    program_type     VARCHAR(20) NOT NULL,
    duration         INTERVAL,
    published_at     TIMESTAMPTZ,
    thumbnail        VARCHAR(2048) NOT NULL DEFAULT '',
    video_url        VARCHAR(2048) NOT NULL DEFAULT '',
    external_id      VARCHAR(100),
    deleted_at       TIMESTAMPTZ,
    status           VARCHAR(15) NOT NULL DEFAULT 'active',
    category_id      BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    language_id      BIGINT REFERENCES languages(id) ON DELETE SET NULL,
    import_source_id BIGINT REFERENCES import_sources(id) ON DELETE SET NULL,
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Partial index for filtering non-deleted programs
CREATE INDEX idx_programs_not_deleted ON programs (id) WHERE deleted_at IS NULL;

CREATE INDEX idx_programs_status_published_at ON programs(status, published_at DESC);
CREATE INDEX idx_programs_external_id ON programs(external_id);
CREATE INDEX idx_programs_category_id ON programs(category_id);
CREATE INDEX idx_programs_language_id ON programs(language_id);
CREATE INDEX idx_programs_import_source_id ON programs(import_source_id);

-- Idempotent import: unique per source+external_id when both present
CREATE UNIQUE INDEX uq_program_import_source_external_id
    ON programs(import_source_id, external_id)
    WHERE import_source_id IS NOT NULL AND external_id IS NOT NULL;

-- Seed programs
INSERT INTO programs (id, title, description, program_type, duration, published_at, thumbnail, video_url, external_id, status, category_id, language_id, import_source_id) VALUES
(
    gen_random_uuid(),
    'التاريخ الظالم لبدو مصر',
    'حلقة بودكاست تناقش التاريخ الظالم لقبائل البدو في مصر.',
    'podcast',
    INTERVAL '1 hour 2 minutes 15 seconds',
    '2026-02-20T10:00:00Z',
    'https://img.youtube.com/vi/placeholder1/maxresdefault.jpg',
    'https://youtu.be/placeholder1',
    'placeholder1',
    'active',
    (SELECT id FROM categories WHERE slug = 'podcast'),
    (SELECT id FROM languages WHERE code = 'ar'),
    (SELECT id FROM import_sources WHERE source_type = 'youtube' LIMIT 1)
),
(
    gen_random_uuid(),
    'زامل السبيعي: همّ القبيلة في صوت شاعر',
    'حلقة بودكاست مع زامل السبيعي عن همّ القبيلة وقضاياها.',
    'podcast',
    INTERVAL '1 hour 39 minutes',
    '2026-02-19T10:00:00Z',
    'https://img.youtube.com/vi/placeholder2/maxresdefault.jpg',
    'https://youtu.be/placeholder2',
    'placeholder2',
    'active',
    (SELECT id FROM categories WHERE slug = 'podcast'),
    (SELECT id FROM languages WHERE code = 'ar'),
    (SELECT id FROM import_sources WHERE source_type = 'youtube' LIMIT 1)
),
(
    gen_random_uuid(),
    'مذكرات محقق حوادث طيران',
    'بودكاست عن التحقيق في حوادث الطيران وقصص واقعية.',
    'podcast',
    INTERVAL '1 hour 15 minutes 30 seconds',
    '2026-02-18T10:00:00Z',
    'https://img.youtube.com/vi/placeholder3/maxresdefault.jpg',
    'https://youtu.be/placeholder3',
    'placeholder3',
    'active',
    (SELECT id FROM categories WHERE slug = 'podcast'),
    (SELECT id FROM languages WHERE code = 'ar'),
    (SELECT id FROM import_sources WHERE source_type = 'youtube' LIMIT 1)
),
(
    gen_random_uuid(),
    'كيف نتغلّب على الشك في مسار حياتنا',
    'حلقة بودكاست تناقش الشك في الحياة واتخاذ القرار.',
    'podcast',
    INTERVAL '59 minutes 46 seconds',
    '2026-02-17T10:00:00Z',
    'https://img.youtube.com/vi/placeholder4/maxresdefault.jpg',
    'https://youtu.be/placeholder4',
    'placeholder4',
    'active',
    (SELECT id FROM categories WHERE slug = 'podcast'),
    (SELECT id FROM languages WHERE code = 'ar'),
    (SELECT id FROM import_sources WHERE source_type = 'youtube' LIMIT 1)
),
(
    gen_random_uuid(),
    'صانع المحتوى الهلالي – عبدالله الشهراني',
    'وثائقي عن صانع المحتوى الهلالي وتأثيره على المشهد الرياضي.',
    'documentary',
    INTERVAL '21 minutes',
    '2026-02-21T10:00:00Z',
    'https://img.youtube.com/vi/placeholder5/maxresdefault.jpg',
    'https://youtu.be/placeholder5',
    'placeholder5',
    'active',
    (SELECT id FROM categories WHERE slug = 'documentary'),
    (SELECT id FROM languages WHERE code = 'ar'),
    (SELECT id FROM import_sources WHERE source_type = 'youtube' LIMIT 1)
);
