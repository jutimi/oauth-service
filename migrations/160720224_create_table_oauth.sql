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

-- +goose Down
DROP TABLE oauth