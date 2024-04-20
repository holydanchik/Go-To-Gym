package main

import (
	"encoding/json"
	"github.com/holydanchik/GoToGym/pkg/go-to-gym/models"
	"net/http"
	"strconv"
)

type WorkoutHandler struct {
	Model *models.WorkoutModel
}

func (wh *WorkoutHandler) GetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса для пагинации, фильтрации и сортировки
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	sortBy := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")
	userID := r.URL.Query().Get("user_id")

	// Получить все записи о тренировках из базы данных с учетом пагинации, фильтрации и сортировки
	workouts, err := wh.Model.GetAll(page, limit, userID, sortBy, sortOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразовать полученные данные в формат JSON и отправить клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}
