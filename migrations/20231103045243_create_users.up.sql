CREATE TABLE users (
    id bigserial not null primary key,
    email varchar not null unique,
    encrypted_password varchar not null
);

CREATE TABLE works (
    id bigserial not null primary key,
    user_id bigint references users (id),
    name varchar not null,
    description varchar,
    created_at timestamp not null,
    expired_at timestamp
);