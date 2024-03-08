package models

import (
	"database/sql"
	"time"
)

type Exercise struct {
	ID        int       `json:"id"`
	WorkoutID int       `json:"workout_id"`
	Name      string    `json:"name"`
	Sets      int       `json:"sets"`
	Reps      int       `json:"reps"`
	Weight    float64   `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
}

type ExerciseModel struct {
	DB *sql.DB
}

func (em *ExerciseModel) Insert(exercise *Exercise) error {
	query := `INSERT INTO exercises (workout_id, name, sets, reps, weight, created_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := em.DB.QueryRow(query, exercise.WorkoutID, exercise.Name, exercise.Sets, exercise.Reps, exercise.Weight, exercise.CreatedAt).Scan(&exercise.ID)
	if err != nil {
		return err
	}
	return nil
}

func (em *ExerciseModel) Get(id int) (*Exercise, error) {
	query := `SELECT id, workout_id, name, sets, reps, weight, created_at FROM exercises WHERE id = $1`
	var exercise Exercise
	err := em.DB.QueryRow(query, id).Scan(&exercise.ID, &exercise.WorkoutID, &exercise.Name, &exercise.Sets, &exercise.Reps, &exercise.Weight, &exercise.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (em *ExerciseModel) Update(exercise *Exercise) error {
	query := `UPDATE exercises SET workout_id = $1, name = $2, sets = $3, reps = $4, weight = $5 WHERE id = $6`
	_, err := em.DB.Exec(query, exercise.WorkoutID, exercise.Name, exercise.Sets, exercise.Reps, exercise.Weight, exercise.ID)
	if err != nil {
		return err
	}
	return nil
}

func (em *ExerciseModel) Delete(id int) error {
	query := `DELETE FROM exercises WHERE id = $1`
	_, err := em.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
