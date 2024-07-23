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
)

-- +goose Down
DROP TABLE oauth