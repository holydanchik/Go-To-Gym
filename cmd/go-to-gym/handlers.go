package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/holydanchik/GoToGym/pkg/go-to-gym/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	Model *models.UserModel
}

// CreateUser обрабатывает запрос на создание нового пользователя.
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Декодируем JSON из тела запроса в структуру User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Применяем хеширование пароля перед сохранением
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.PasswordHash = hashedPassword

	// Устанавливаем текущее время как время создания пользователя
	user.CreatedAt = time.Now()

	// Вставляем пользователя в базу данных
	err = uh.Model.Insert(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ с данными пользователя в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// hashPassword хеширует пароль с использованием bcrypt.
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uh.Model.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = id

	err = uh.Model.Update(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uh.Model.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса для пагинации, фильтрации и сортировки
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	filter := r.URL.Query().Get("filter")
	sortBy := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")

	// Получить все записи о пользователях из базы данных
	users, err := uh.Model.GetAll(page, limit, filter, sortBy, sortOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразовать полученные данные в формат JSON и отправить клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
