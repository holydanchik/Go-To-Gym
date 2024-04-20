package models

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) Insert(user *User) error {
	query := `INSERT INTO users (username, email, password_hash, created_at)
              VALUES ($1, $2, $3, $4) RETURNING id`
	err := um.DB.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserModel) Get(id int) (*User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1`
	var user User
	err := um.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (um *UserModel) Update(user *User) error {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3`
	_, err := um.DB.Exec(query, user.Username, user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserModel) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := um.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserModel) GetAll(page, limit int, filter, sortBy, sortOrder string) ([]*User, error) {
	// Формируем SQL запрос с учетом параметров пагинации, фильтрации и сортировки
	query := "SELECT id, username, email, created_at FROM users"

	if filter != "" {
		query += fmt.Sprintf(" WHERE username LIKE '%%%s%%' OR email LIKE '%%%s%%'", filter, filter)
	}

	if sortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	}

	if limit > 0 {
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	// Выполняем запрос к базе данных
	rows, err := um.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Обрабатываем результат запроса и возвращаем список пользователей
	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
