package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/holydanchik/GoToGym/pkg/go-to-gym/models"
	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config       config
	userModel    *models.UserModel
	workoutModel *models.WorkoutModel
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://danchik@localhost/gym?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize UserModel
	userModel := &models.UserModel{DB: db}
	workoutModel := &models.WorkoutModel{DB: db}

	app := &application{
		config:       cfg,
		userModel:    userModel,
		workoutModel: workoutModel,
	}

	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// User Singleton
	userHandler := &UserHandler{Model: app.userModel}
	workoutHandler := &WorkoutHandler{Model: app.workoutModel}

	// Create a new user
	v1.HandleFunc("/users/register", userHandler.CreateUser).Methods("POST")
	v1.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	// Get a specific user
	v1.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET")
	// Update a specific user
	v1.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	// Delete a specific user
	v1.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
	// New route for workouts
	v1.HandleFunc("/workouts", workoutHandler.GetAllWorkouts).Methods("GET")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
