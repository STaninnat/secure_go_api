-- +goose Up
CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        name TEXT NOT NULL,
        password VARCHAR(255) NOT NULL,
        api_key VARCHAR(64) UNIQUE NOT NULL,
        expires_at TIMESTAMP NOT NULL
    );

-- +goose Down
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;