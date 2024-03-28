CREATE TABLE IF NOT EXISTS goals
(
    id      INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    user_id INT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name    VARCHAR(35) NOT NULL,
    time_seconds BIGINT NOT NULL,
    date_start timestamp NOT NULL,
    date_end   timestamp NOT NULL
)