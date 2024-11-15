-- +goose Up
CREATE TABLE posts (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP  NOT NULL,
    updated_at TIMESTAMP  NOT NULL,
    post TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS posts;