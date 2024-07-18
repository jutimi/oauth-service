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

-- +goose Down
DROP TABLE users;