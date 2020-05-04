
-- +migrate Up
CREATE TABLE IF NOT EXISTS songs (
--     id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL UNSIGNED,
--     id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    id BIGINT AUTO_INCREMENT UNSIGNED NOT NULL,
    title varchar(255) NOT NULL,
    artist varchar(255) NOT NULL,
    music_age int NOT NULL,
    image varchar(255),
    video varchar(255),
    album varchar(255),
    description varchar(255),
    spotify_id BIGINT,
    user_id BIGINT UNSIGNED NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    deleted_at timestamp,
    PRIMARY KEY (id),
    FOREIGN KEY(user_id) -- 外部キー
    REFERENCES users(id)
);
-- 以下は後から追加する場合
-- ALTER TABLE songs ADD FOREIGN KEY user_id REFERENCES users(id);

-- +migrate Down
DROP TABLE IF EXISTS songs;
