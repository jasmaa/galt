CREATE DATABASE galt;

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    profile_img_url VARCHAR(255) NOT NULL
);

CREATE TABLE statuses (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    content VARCHAR(255) NOT NULL,
    posted_timestamp TIMESTAMP NOT NULL,
    is_edited BOOLEAN NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE comments (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    status_id UUID NOT NULL,
    parent_comment_id UUID,
    content VARCHAR(255) NOT NULL, 
    posted_timestamp TIMESTAMP NOT NULL,
    is_edited BOOLEAN NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(status_id) REFERENCES statuses(id),
    FOREIGN KEY(parent_comment_id) REFERENCES comments(id)
);

CREATE TABLE circles (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE circle_user_pairs (
    user_id UUID NOT NULL,
    circle_id UUID NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(circle_id) REFERENCES circles(id) ON DELETE CASCADE
);

CREATE TABLE status_like_pairs (
    user_id UUID NOT NULL,
    status_id UUID NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(status_id) REFERENCES statuses(id) ON DELETE CASCADE
);

CREATE TABLE comment_like_pairs (
    user_id UUID NOT NULL,
    comment_id UUID NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE
);