CREATE TYPE role_type AS ENUM ('user', 'admin');

CREATE TABLE IF NOT EXISTS users
(
    id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name     VARCHAR(35)  NOT NULL,
    email    VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL
);


CREATE TABLE IF NOT EXISTS projects
(
    id      INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name    VARCHAR(35) NOT NULL
);

CREATE TABLE IF NOT EXISTS entries
(
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    INT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    project_id INT       NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name       TEXT DEFAULT '',
    time_start TIMESTAMP NOT NULL,
    time_end   TIMESTAMP NOT NULL
);

insert into users (name, email, password)
values ('test', 'test', 'password');
