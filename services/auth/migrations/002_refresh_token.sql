CREATE TABLE referesh_tokens (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_At TIMESTAMP NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMP,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE INDEX idx_referesh_tokens_user_id ON referesh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON referesh_tokens(token);