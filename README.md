# Go To Gym
by Daniyal Tuzelbayev 21B030935

# How to run an app directly on Terminal
Provide all needed correct values.
```
go run ./cmd/go-to-gym \
-dsn="postgres://username@localhost/gym?sslmode=disable" \
-migrations=file://pkg/go-to-gym/migrations \
-env=development \
-port=4000
```

## Introduction
Go To Gym is a fitness application designed to help users plan and track their workouts effectively. With Go To Gym, users can create personalized training programs, log their workouts, track their progress to monitor their achievements.

## DB Structure
``` sql
CREATE TABLE IF NOT EXISTS workouts
(
    id              SERIAL PRIMARY KEY,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    exercises       TEXT[] NOT NULL,
    calories_burned INTEGER,
    version         INTEGER NOT NULL DEFAULT 1
);
CREATE TABLE IF NOT EXISTS exercises
(
    id         SERIAL PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name       VARCHAR(255)                NOT NULL,
    sets       INT                         NOT NULL,
    reps       INT                         NOT NULL,
    version    integer                     NOT NULL DEFAULT 1,
    workout_id INT REFERENCES workouts (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS users
(
    id            BIGSERIAL PRIMARY KEY,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name          TEXT NOT NULL,
    email         TEXT UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    activated     BOOLEAN NOT NULL,
    version       INTEGER NOT NULL DEFAULT 1
);
```

## Add write workouts permission to user example SQL
``` sql
INSERT INTO users_permissions
VALUES (
           (SELECT id FROM users WHERE email = 'email@example.com'),
           (SELECT id FROM permissions WHERE code = 'workouts:write')
       );

```

# API Endpoints
## Workouts
```
GET /v1/workouts: Retrieve all workouts.
POST /v1/workouts: Create a new workout.
GET /v1/workouts/{id}: Retrieve a specific workout by ID.
PATCH /v1/workouts/{id}: Update an existing workout.
DELETE /v1/workouts/{id}: Delete a workout.
```
## Exercises
```
GET /v1/exercises: Retrieve all exercises.
POST /v1/exercises: Create a new exercise.
GET /v1/exercises/{id}: Retrieve a specific exercise by ID.
PATCH /v1/exercises/{id}: Update an existing exercise.
DELETE /v1/exercise/{id}: Delete an exercise.
GET /v1/workouts/{id}/exercises: Retrieve all exercises that attached to specific workout_id.
```
## Users
```
POST /v1/users: Register a new user.
PUT /v1/users/activated: Activate a user.
```
## Authentication
```
POST /v1/tokens/authentication: Create an authentication token.
```
## Authorization
Each API endpoint is guarded by specific permissions:
```
workouts:read: Read permission for workouts.
workouts:write: Write permission for workouts.
These permissions are enforced using the requirePermission middleware.
```

## Contributing
Contributions to Go To Gym are welcome! Feel free to open issues for bug fixes, feature requests, or any other improvements you'd like to see. Pull requests are also encouraged.




