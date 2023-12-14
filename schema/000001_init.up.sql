CREATE TABLE users
(
    id serial not null unique,
    name varchar(100) not null,
    surname varchar(100),
    email varchar(255) not null unique,
    password_hash varchar(255) not null,
    is_confirmed boolean default FALSE
);

CREATE TABLE posts
(
    id serial not null unique,
    title varchar(255),
    content text not null,
    users_id int references users (id) on delete set null,
    date timestamp default CURRENT_TIMESTAMP
);

CREATE TABLE comments
(
    id serial not null unique,
    content text not null,
    users_id int references users (id) on delete set null,
    posts_id int references posts (id) on delete cascade,
    date timestamp default CURRENT_TIMESTAMP
);

CREATE TABLE likes_posts
(
    id serial not null unique,
    users_id int references users (id) on delete cascade,
    posts_id int references posts (id) on delete cascade,
    date timestamp default CURRENT_TIMESTAMP
);

CREATE TABLE auth_emails
(
    id serial not null unique,
    users_id int references users (id) on delete cascade,
    hash varchar(255)
)