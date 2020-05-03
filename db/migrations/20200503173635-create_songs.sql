
-- +migrate Up
CREATE TABLE IF NOT EXISTS songs (
--     id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL UNSIGNED,
--     id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    id BIGINT AUTO_INCREMENT NOT NULL,
    title varchar(255) NOT NULL,
    artist varchar(255) NOT NULL,
    music_age int NOT NULL,
    image varchar(255),
    video varchar(255),
    album varchar(255),
    description varchar(255),
    spotify_id BIGINT,
    user_id BIGINT NOT NULL,
--     user_id BIGINT FOREIGN KEY ('id') REFERENCES 'users'


--     外部キー設定したい
--     user_id int,
--     user_id BIGINT FOREIGN KEY ('id') REFERENCES 'users',
--     CONSTRAINT 'songs_user_id_users_id_foreign' FOREIGN KEY ('user_id') REFERENCES 'users' ('id')
--     email varchar(255),
--     date timestamp NOT NULL,
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
