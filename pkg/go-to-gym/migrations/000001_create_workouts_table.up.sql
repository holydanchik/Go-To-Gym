CREATE TABLE IF NOT EXISTS workouts
(
    id              SERIAL PRIMARY KEY,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name            VARCHAR(255)                NOT NULL,
    description     TEXT,
    exercises       TEXT[]                      not null,
    calories_burned INTEGER,
    version         integer                     NOT NULL DEFAULT 1
);