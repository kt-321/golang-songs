
-- +migrate Up
CREATE TABLE IF NOT EXISTS songs (
    id BIGINT AUTO_INCREMENT NOT NULL,
    title varchar(255) NOT NULL,
    artist varchar(255) NOT NULL,
    music_age int NOT NULL,
    image varchar(255),
    video varchar(255),
    album varchar(255),
    description varchar(255),
    spotify_track_id varchar(255),
    user_id BIGINT NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    deleted_at timestamp,
    PRIMARY KEY (id),
    FOREIGN KEY(user_id)
    REFERENCES users(id)
);
-- +migrate Down
DROP TABLE IF EXISTS songs;