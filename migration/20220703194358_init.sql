-- +goose Up
CREATE TABLE users
(
    id            SERIAL NOT NULL UNIQUE,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY (id)
);

CREATE TABLE posts
(
    id           SERIAL NOT NULL UNIQUE,
    user_id      INTEGER NOT NULL,
    reading_time INTEGER,
    status       INTEGER,
    title        VARCHAR(255) NOT NULL,
    subtitle     TEXT,
    image_url    TEXT,
    content      TEXT,
    slug         VARCHAR(255),
    published_at TIMESTAMP,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL,
    deleted_at   TIMESTAMP,
    CONSTRAINT pk_posts PRIMARY KEY (id)
);

CREATE TABLE tags
(
    id    SERIAL NOT NULL UNIQUE,
    name  VARCHAR(255) NOT NULL UNIQUE,
    slug  VARCHAR(255) NOT NULL,
    CONSTRAINT pk_tags PRIMARY KEY (id)
);

CREATE TABLE posts_tags
(
    id      SERIAL NOT NULL UNIQUE,
    tag_id  INTEGER REFERENCES tags (id) ON DELETE CASCADE NOT NULL,
    post_id INTEGER REFERENCES posts (id) ON DELETE CASCADE NOT NULL,
    CONSTRAINT pk_posts_tags PRIMARY KEY (id)
);
-- +goose Down
DROP TABLE posts_tags;
DROP TABLE posts;
DROP TABLE users;
DROP TABLE tags;