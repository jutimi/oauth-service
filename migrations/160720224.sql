CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email varchar(100),
    phone_number varchar(20),
    password text NOT NULL,
    is_active boolean DEFAULT true
    created_at bigint,
    updated_at bigint
);

CREATE TABLE oauth (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    ip text NOT NULL,
    platform text NOT NULL,
    token text NOT NULL,
    status text NOT NULL,
    expire_at bigint NOT NULL,
    login_at bigint,
    created_at bigint,
    updated_at bigint
)