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

-- +goose Down
DROP TABLE users;