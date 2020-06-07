CREATE DATABASE galt;

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    profile_img_url VARCHAR(255) NOT NULL
);

CREATE TABLE statuses (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    content VARCHAR(255) NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE comments (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    status_id VARCHAR(255) NOT NULL,
    parent_comment_id VARCHAR(255) NOT NULL,
    content VARCHAR(255) NOT NULL, 
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(status_id) REFERENCES statuses(id),
    FOREIGN KEY(parent_comment_id) REFERENCES comments(id)
);

CREATE TABLE circles (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE circle_user_pairs (
    user_id VARCHAR(255) NOT NULL,
    circle_id VARCHAR(255) NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(circle_id) REFERENCES circles(id)
);