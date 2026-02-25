CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Seed default roles
INSERT INTO roles (id, name) VALUES
    (gen_random_uuid(), 'admin'),
    (gen_random_uuid(), 'editor'),
    (gen_random_uuid(), 'user');

-- Seed default permissions
INSERT INTO permissions (id, code) VALUES
    (gen_random_uuid(), 'content:create'),
    (gen_random_uuid(), 'content:update'),
    (gen_random_uuid(), 'content:delete');

-- Assign all permissions to admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r CROSS JOIN permissions p WHERE r.name = 'admin';

-- Assign content:create and content:update to editor
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'editor' AND p.code IN ('content:create', 'content:update');

-- Seed admin user (password: Admin@1234)
INSERT INTO users (id, email, password_hash, status) VALUES
(gen_random_uuid(), 'admin@test.com', '$2a$12$eM11HfXzop1zNjCXpFnWXe/w.O6DAqfN8MPr3tTODoxUohOf8scHy', 'active');

-- Assign admin role to admin user
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r WHERE u.email = 'admin@test.com' AND r.name = 'admin';
