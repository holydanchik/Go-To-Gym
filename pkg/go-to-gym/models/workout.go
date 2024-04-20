package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Workout struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Duration  int       `json:"duration"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkoutModel struct {
	DB *sql.DB
}

func (wm *WorkoutModel) Insert(workout *Workout) error {
	query := `INSERT INTO workouts (user_id, name, duration, date, created_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := wm.DB.QueryRow(query, workout.UserID, workout.Name, workout.Duration, workout.Date, workout.CreatedAt).Scan(&workout.ID)
	if err != nil {
		return err
	}
	return nil
}

func (wm *WorkoutModel) Get(id int) (*Workout, error) {
	query := `SELECT id, user_id, name, duration, date, created_at FROM workouts WHERE id = $1`
	var workout Workout
	err := wm.DB.QueryRow(query, id).Scan(&workout.ID, &workout.UserID, &workout.Name, &workout.Duration, &workout.Date, &workout.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (wm *WorkoutModel) Update(workout *Workout) error {
	query := `UPDATE workouts SET user_id = $1, name = $2, duration = $3, date = $4 WHERE id = $5`
	_, err := wm.DB.Exec(query, workout.UserID, workout.Name, workout.Duration, workout.Date, workout.ID)
	if err != nil {
		return err
	}
	return nil
}

func (wm *WorkoutModel) Delete(id int) error {
	query := `DELETE FROM workouts WHERE id = $1`
	_, err := wm.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (wm *WorkoutModel) GetAll(page, limit int, userID, sortBy, sortOrder string) ([]*Workout, error) {
	query := "SELECT id, user_id, name, duration, date, created_at FROM workouts"

	// Добавляем условие WHERE для фильтрации по user_id
	if userID != "" {
		query += fmt.Sprintf(" WHERE user_id = %s", userID)
	}

	// Добавляем сортировку
	if sortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	}

	// Добавляем пагинацию
	if limit > 0 {
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := wm.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []*Workout
	for rows.Next() {
		var workout Workout
		err := rows.Scan(&workout.ID, &workout.UserID, &workout.Name, &workout.Duration, &workout.Date, &workout.CreatedAt)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, &workout)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}
