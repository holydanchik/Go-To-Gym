package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const PORT = ":8000"

func main() {
	log.Println("starting API server")
	router := mux.NewRouter()
	http.Handle("/", router)

	log.Println("Server started on port", PORT)
	err := http.ListenAndServe(PORT, router)
	log.Fatal(err)
}
