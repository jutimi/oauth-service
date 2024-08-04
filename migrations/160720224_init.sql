CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(100) UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    password TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    limit_workspace INTEGER DEFAULT 1,
    created_at BIGINT,
    updated_at BIGINT
);

CREATE TABLE permissions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    workspace_id uuid,
    user_workspace_id uuid,
    scopes text NOT NULL,
    created_at bigint,
    updated_at bigint
);

CREATE TABLE oauth (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    ip VARCHAR(10) NOT NULL,
    platform VARCHAR(20) NOT NULL,
    token TEXT NOT NULL,
    status VARCHAR(10) NOT NULL,
    expire_at BIGINT NOT NULL,
    login_at BIGINT,
    scope VARCHAR(10) NOT NULL,
    created_at BIGINT,
    updated_at BIGINT
);