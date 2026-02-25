CREATE TABLE languages (
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    code VARCHAR(10) NOT NULL UNIQUE
);

-- Seed languages
INSERT INTO languages (name, code) VALUES
    ('العربية', 'ar'),
    ('English', 'en');
