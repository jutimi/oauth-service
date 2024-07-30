CREATE TABLE permissions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    workspace_id uuid,
    user_workspace_id uuid,
    scopes text NOT NULL,
    created_at bigint,
    updated_at bigint
);

DROP TABLE permissions;