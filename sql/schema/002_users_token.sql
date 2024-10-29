-- +goose Up
CREATE TABLE 
    users_token (
        id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        access_token_expires_at TIMESTAMP NOT NULL,
        refresh_token VARCHAR(255) UNIQUE NOT NULL,
        refresh_token_expires_at TIMESTAMP NOT NULL,
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE IF EXISTS users_token;