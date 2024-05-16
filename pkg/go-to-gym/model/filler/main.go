package filler

import "github.com/holydanchik/GoToGym/pkg/go-to-gym/model"

func PopulateDatabase(models model.Models) error {
	for _, workout := range workouts {
		err := models.Workouts.Insert(&workout)
		if err != nil {
			return err
		}
	}

	for _, exercise := range exercises {
		err := models.Exercises.Insert(&exercise)
		if err != nil {
			return err
		}
	}
	return nil
}

var workouts = []model.Workout{
	{Name: "Legs", Description: "Legs + Arms program", Exercises: []string{"Squats", "Lunges", "Leg Press", "Bicep Curls", "Dips", "Shoulder Press"}, CaloriesBurned: 520},
	{Name: "Chest", Description: "Chest + Core program", Exercises: []string{"Bench Press", "Push-ups", "Dumbbell Flyes", "Planks", "Russian Twists", "Leg Raises"}, CaloriesBurned: 400},
	{Name: "Back", Description: "Back Day program", Exercises: []string{"Deadlifts", "Pull-ups", "Rows"}, CaloriesBurned: 250},
	{Name: "Cardio", Description: "Cardio Workout program", Exercises: []string{"Running", "Cycling", "Jumping Jacks"}, CaloriesBurned: 300},
	{Name: "Full Body", Description: "Full Body Workout program", Exercises: []string{"Squats", "Push-ups", "Pull-ups", "Planks"}, CaloriesBurned: 350},
}

var exercises = []model.Exercise{
	{Name: "Squats", Sets: 3, Reps: 5, WorkoutID: 1},
	{Name: "Lunges", Sets: 2, Reps: 12, WorkoutID: 1},
	{Name: "Leg Press", Sets: 4, Reps: 12, WorkoutID: 1},
	{Name: "Bicep Curls", Sets: 3, Reps: 12, WorkoutID: 1},
	{Name: "Dips", Sets: 3, Reps: 10, WorkoutID: 1},
	{Name: "Shoulder Press", Sets: 2, Reps: 20, WorkoutID: 1},
	{Name: "Bench Press", Sets: 4, Reps: 8, WorkoutID: 2},
	{Name: "Push-ups", Sets: 3, Reps: 15, WorkoutID: 2},
	{Name: "Dumbbell Flyes", Sets: 3, Reps: 12, WorkoutID: 2},
	{Name: "Planks", Sets: 3, Reps: 60, WorkoutID: 2},
	{Name: "Russian Twists", Sets: 3, Reps: 20, WorkoutID: 2},
	{Name: "Leg Raises", Sets: 3, Reps: 15, WorkoutID: 2},
	{Name: "Deadlifts", Sets: 4, Reps: 6, WorkoutID: 3},
	{Name: "Pull-ups", Sets: 3, Reps: 10, WorkoutID: 3},
	{Name: "Rows", Sets: 3, Reps: 12, WorkoutID: 3},
	{Name: "Running", Sets: 1, Reps: 30, WorkoutID: 4},
	{Name: "Cycling", Sets: 1, Reps: 30, WorkoutID: 4},
	{Name: "Jumping Jacks", Sets: 1, Reps: 60, WorkoutID: 4},
	{Name: "Squats", Sets: 3, Reps: 10, WorkoutID: 5},
	{Name: "Push-ups", Sets: 3, Reps: 20, WorkoutID: 5},
	{Name: "Pull-ups", Sets: 3, Reps: 10, WorkoutID: 5},
	{Name: "Planks", Sets: 3, Reps: 60, WorkoutID: 5},
}
