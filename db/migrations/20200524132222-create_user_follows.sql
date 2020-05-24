
-- +migrate Up
CREATE TABLE IF NOT EXISTS user_follows (
    id BIGINT AUTO_INCREMENT NOT NULL,
    user_id BIGINT NOT NULL,
    follow_id BIGINT NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    deleted_at timestamp,
    PRIMARY KEY (id),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(follow_id) REFERENCES users(id)
);
-- +migrate Down
DROP TABLE IF EXISTS user_follows;