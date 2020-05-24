
-- +migrate Up
CREATE TABLE IF NOT EXISTS bookmarks (
    id BIGINT AUTO_INCREMENT NOT NULL,
    user_id BIGINT NOT NULL,
    song_id BIGINT NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    deleted_at timestamp,
    PRIMARY KEY (id),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(song_id) REFERENCES songs(id)
);
-- +migrate Down
DROP TABLE IF EXISTS bookmarks;