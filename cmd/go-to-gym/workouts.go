package main

import (
	"errors"
	"fmt"
	"github.com/holydanchik/GoToGym/pkg/go-to-gym/model"
	"github.com/holydanchik/GoToGym/pkg/go-to-gym/validator"
	"net/http"
)

func (app *application) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string   `json:"name"`
		Description    string   `json:"description,omitempty"`
		Exercises      []string `json:"exercises"`
		CaloriesBurned int      `json:"calories_burned,omitempty"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	workout := &model.Workout{
		Name:           input.Name,
		Description:    input.Description,
		Exercises:      input.Exercises,
		CaloriesBurned: input.CaloriesBurned,
	}
	v := validator.New()
	if model.ValidateWorkout(v, workout); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	if err := app.models.Workouts.Insert(workout); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/workouts/%d", workout.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"workout": workout}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	workout, err := app.models.Workouts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"workout": workout}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	workout, err := app.models.Workouts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Name           *string   `json:"name"`
		Description    *string   `json:"description,omitempty"`
		Exercises      *[]string `json:"exercises,omitempty"`
		CaloriesBurned *int      `json:"calories_burned,omitempty"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Name != nil {
		workout.Name = *input.Name
	}
	if input.Description != nil {
		workout.Description = *input.Description
	}
	if input.Exercises != nil {
		workout.Exercises = *input.Exercises
	}
	if input.CaloriesBurned != nil {
		workout.CaloriesBurned = *input.CaloriesBurned
	}
	v := validator.New()
	if model.ValidateWorkout(v, workout); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	if err := app.models.Workouts.Update(workout); err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"workout": workout}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Workouts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "workout successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string
		Filters   model.Filters
		Exercises string
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "calories_burned", "-id", "-name", "-calories_burned"}
	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	input.Exercises = app.readString(qs, "exercises", "")
	workouts, metadata, err := app.models.Workouts.GetAll(input.Name, input.Filters, input.Exercises) //input.Exercises)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"workouts": workouts, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
