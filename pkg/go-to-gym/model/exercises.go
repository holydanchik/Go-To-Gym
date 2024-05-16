package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/holydanchik/GoToGym/pkg/go-to-gym/validator"
	"time"
)

type Exercise struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Sets      int       `json:"sets"`
	Reps      int       `json:"reps"`
	WorkoutID int       `json:"workout_id,omitempty"`
	Version   int       `json:"version"`
}

func ValidateExercise(v *validator.Validator, e *Exercise) {
	v.Check(e.Name != "", "name", "must be provided")
	v.Check(len(e.Name) <= 100, "name", "must not be more than 100 characters long")
	v.Check(e.Sets >= 0, "sets", "must be a non-negative value")
	v.Check(e.Reps >= 0, "reps", "must be a non-negative value")
}

type ExerciseModel struct {
	DB *sql.DB
}

func (m ExerciseModel) Insert(exercise *Exercise) error {
	query := `
		INSERT INTO exercises (name, sets, reps, workout_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{exercise.Name, exercise.Sets, exercise.Reps, exercise.WorkoutID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&exercise.ID, &exercise.CreatedAt, &exercise.Version)
}

func (m ExerciseModel) Get(id int64) (*Exercise, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, sets, reps, workout_id, version
		FROM exercises
		WHERE id = $1`

	var exercise Exercise

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.CreatedAt,
		&exercise.Name,
		&exercise.Sets,
		&exercise.Reps,
		&exercise.WorkoutID,
		&exercise.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &exercise, nil
}

func (m ExerciseModel) Update(exercise *Exercise) error {
	query := `
		UPDATE exercises
		SET name = $1, sets = $2, reps = $3, workout_id = $4, version = version + 1
		WHERE id = $5 and version = $6
		RETURNING version`

	args := []interface{}{
		exercise.Name,
		exercise.Sets,
		exercise.Reps,
		exercise.WorkoutID,
		exercise.ID,
		exercise.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&exercise.Version)
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

func (m ExerciseModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM exercises
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

func (m ExerciseModel) GetAll(name string, paramWorkoutID int, from, to int, filters Filters) ([]*Exercise, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, name, sets, reps, workout_id, version
		FROM exercises
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND workout_id = $2
		AND (sets >= $3 OR $3 = 0)
		AND (sets <= $4 OR $4 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, paramWorkoutID, from, to, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var exercises []*Exercise

	for rows.Next() {
		var exercise Exercise

		err := rows.Scan(
			&totalRecords,
			&exercise.ID,
			&exercise.CreatedAt,
			&exercise.Name,
			&exercise.Sets,
			&exercise.Reps,
			&exercise.WorkoutID,
			&exercise.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		exercises = append(exercises, &exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return exercises, metadata, nil
}
