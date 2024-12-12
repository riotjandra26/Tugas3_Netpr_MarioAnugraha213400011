package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Film struct
type Film struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

// Slice untuk menyimpan data film
var films []Film
var idCounter = 1

// Get all films
func getFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(films)
}

// Get single film by ID
func getFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, film := range films {
		if film.ID == id {
			json.NewEncoder(w).Encode(film)
			return
		}
	}
	http.Error(w, "Film not found", http.StatusNotFound)
}

// Create a new film
func createFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var film Film
	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	film.ID = idCounter
	idCounter++
	films = append(films, film)
	json.NewEncoder(w).Encode(film)
}

// Update film by ID
func updateFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for index, film := range films {
		if film.ID == id {
			films = append(films[:index], films[index+1:]...)
			var updatedFilm Film
			err := json.NewDecoder(r.Body).Decode(&updatedFilm)
			if err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			updatedFilm.ID = id
			films = append(films, updatedFilm)
			json.NewEncoder(w).Encode(updatedFilm)
			return
		}
	}
	http.Error(w, "Film not found", http.StatusNotFound)
}

// Delete film by ID
func deleteFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for index, film := range films {
		if film.ID == id {
			films = append(films[:index], films[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Film not found", http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/films", getFilms).Methods("GET")
	r.HandleFunc("/films/{id}", getFilm).Methods("GET")
	r.HandleFunc("/films", createFilm).Methods("POST")
	r.HandleFunc("/films/{id}", updateFilm).Methods("PUT")
	r.HandleFunc("/films/{id}", deleteFilm).Methods("DELETE")

	// Start server
	http.ListenAndServe(":8080", r)
}
