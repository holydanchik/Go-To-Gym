package main

import (
	"errors"
	"fmt"
	"github.com/holydanchik/GoToGym/pkg/go-to-gym/model"
	"github.com/holydanchik/GoToGym/pkg/go-to-gym/validator"
	"net/http"
)

func (app *application) createExerciseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string `json:"name"`
		Sets      int    `json:"sets"`
		Reps      int    `json:"reps"`
		WorkoutID int    `json:"workout_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	exercise := &model.Exercise{
		Name:      input.Name,
		Sets:      input.Sets,
		Reps:      input.Reps,
		WorkoutID: input.WorkoutID,
	}

	v := validator.New()
	if model.ValidateExercise(v, exercise); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Exercises.Insert(exercise)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/exercises/%d", exercise.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"exercise": exercise}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showExerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	exercise, err := app.models.Exercises.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"exercise": exercise}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateExerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	exercise, err := app.models.Exercises.Get(id)
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
		Name      *string `json:"name"`
		Sets      *int    `json:"sets"`
		Reps      *int    `json:"reps"`
		WorkoutID *int    `json:"workout_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		exercise.Name = *input.Name
	}
	if input.Sets != nil {
		exercise.Sets = *input.Sets
	}
	if input.Reps != nil {
		exercise.Reps = *input.Reps
	}
	if input.WorkoutID != nil {
		exercise.WorkoutID = *input.WorkoutID
	}

	v := validator.New()
	if model.ValidateExercise(v, exercise); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Exercises.Update(exercise)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"exercise": exercise}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteExerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Exercises.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "exercise successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listExercisesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name    string
		SetFrom int
		SetTo   int
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.SetFrom = app.readInt(qs, "setFrom", 0, v)
	input.SetTo = app.readInt(qs, "setTo", 0, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "sets", "reps", "-id", "-name", "-sets", "-reps"}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	exercises, metadata, err := app.models.Exercises.GetAll(input.Name, int(id), input.SetFrom, input.SetTo, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workouts": exercises, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
