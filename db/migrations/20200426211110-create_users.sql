
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    name varchar(255),
    email varchar(255) NOT NULL,
    age int,
    gender int,
    image_url varchar(255),
    favorite_music_age int,
    favorite_artist varchar(255),
    comment varchar(255),
    password varchar(255) NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    deleted_at timestamp
);
-- +migrate Down
DROP TABLE IF EXISTS users;