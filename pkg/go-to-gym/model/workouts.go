package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/holydanchik/GoToGym/pkg/go-to-gym/validator"
	"github.com/lib/pq"
	"time"
)

type Workout struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"-"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	Exercises      []string  `json:"exercises"`
	CaloriesBurned int       `json:"calories_burned,omitempty"`
	Version        int       `json:"version"`
}

func ValidateWorkout(v *validator.Validator, w *Workout) {
	v.Check(w.Name != "", "name", "must be provided")
	v.Check(len(w.Name) <= 100, "name", "must not be more than 100 characters long")
	v.Check(len(w.Exercises) > 0, "exercises", "at least one exercise must be provided")
	v.Check(w.CaloriesBurned >= 0, "calories_burned", "must be a non-negative value")
}

type WorkoutModel struct {
	DB *sql.DB
}

func (m WorkoutModel) Insert(workout *Workout) error {
	query := `
		INSERT INTO workouts (name, description, exercises, calories_burned)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{workout.Name, workout.Description, pq.Array(workout.Exercises), workout.CaloriesBurned}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&workout.ID, &workout.CreatedAt, &workout.Version)
}

func (m WorkoutModel) Get(id int64) (*Workout, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, created_at, name, description, exercises, calories_burned, version
		FROM workouts
		WHERE id = $1`

	var workout Workout

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&workout.ID,
		&workout.CreatedAt,
		&workout.Name,
		&workout.Description,
		pq.Array(&workout.Exercises),
		&workout.CaloriesBurned,
		&workout.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &workout, nil
}

func (m WorkoutModel) Update(workout *Workout) error {
	query := `
		UPDATE workouts
		SET name = $1, description = $2, exercises = $3, calories_burned = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []interface{}{
		workout.Name,
		workout.Description,
		pq.Array(workout.Exercises),
		workout.CaloriesBurned,
		workout.ID,
		workout.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&workout.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m WorkoutModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM workouts
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m WorkoutModel) GetAll(name string, exercises []string, from, to int, filters Filters) ([]*Workout, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, name, description, exercises, calories_burned, version
		FROM workouts
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (exercises @> $2 OR $2 = '{}')
		AND (calories_burned >= $3 OR $3 = 0)
		AND (calories_burned <= $4 OR $4 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, pq.Array(exercises), from, to, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	var workouts []*Workout

	for rows.Next() {
		var workout Workout

		err := rows.Scan(
			&totalRecords,
			&workout.ID,
			&workout.CreatedAt,
			&workout.Name,
			&workout.Description,
			pq.Array(&workout.Exercises),
			&workout.CaloriesBurned,
			&workout.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		workouts = append(workouts, &workout)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return workouts, metadata, nil
}
