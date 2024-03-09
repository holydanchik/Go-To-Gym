# Go To Gym
by Daniyal Tuzelbayev 21B030935

## Introduction
Go To Gym is a fitness application designed to help users plan and track their workouts effectively. With Go To Gym, users can create personalized training programs, log their workouts, track their progress to monitor their achievements.

## Users REST API
Users REST API

POST /api/v1/users         - Create a new user
GET /api/v1/users/{id}     - Get information about a specific user
PUT /api/v1/users/{id}     - Update information about a specific user
DELETE /api/v1/users/{id}  - Delete a specific user


## DB Structure
// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs

Table users {
  id serial [primary key]
  username varchar(50) [not null, unique]
  email varchar(100) [not null, unique]
  password_hash varchar(100) [not null]
  created_at timestamp [default: `CURRENT_TIMESTAMP`]
}

Table workouts {
  id serial [primary key]
  user_id int [not null]
  name varchar(255) [not null]
  duration int [not null]
  date date [not null]
  created_at timestamp [default: `CURRENT_TIMESTAMP`]
}

Table exercises {
  id serial [primary key]
  workout_id int [not null]
  name varchar(255) [not null]
  sets int [not null]
  reps int [not null]
  weight decimal(10,2)
  created_at timestamp [default: `CURRENT_TIMESTAMP`]
}

Ref: workouts.user_id > users.id
Ref: exercises.workout_id > workouts.id


## Features
- User Authentication: Secure registration and authentication system to manage user accounts.
- Workout Management: Create, view, edit, and delete personalized workout programs.
- Workout Logging: Log completed workouts, including exercise details, duration, and intensity.
- Progress Tracking: Track progress over time, including weightlifting progress, cardio performance, and overall fitness improvement.
- Goal Setting: Set fitness goals and milestones, and monitor progress towards achieving them.
- Reporting: Generate detailed reports and statistics to visualize workout data and progress.

## Technologies Used
- **Go (Golang)**: Backend development language for building robust and efficient server-side applications.
- **PostgreSQL**: Relational database management system for storing user data, workout programs, and progress records.
- **RESTful API**: API endpoints implemented using Go's standard library for handling HTTP requests and responses.
- **JWT Authentication**: JSON Web Tokens used for user authentication and authorization.
- **Docker**: Containerization tool for packaging the application and its dependencies into portable containers.
- **Standard Layout**: Project follows the standard layout conventions for organizing Go projects.

## Getting Started
To get started with Go To Gym, follow these steps:

1. Clone the repository: git clone https://github.com/holydanchik/GoToGym.git 
2. Install Docker on your machine if you haven't already.
3. Navigate to the project directory and run the following command to build and run the Docker container: docker-compose up --build
4. Once the container is running, you can access the API at `http://localhost:8000`.
5. Use tools like Postman or curl to interact with the API endpoints and start planning and tracking your workouts!

## Contributing
Contributions to Go To Gym are welcome! Feel free to open issues for bug fixes, feature requests, or any other improvements you'd like to see. Pull requests are also encouraged.




