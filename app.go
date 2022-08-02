package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Sets up database connection, creates a new gorilla/mux router, and sets up routes
func (app *App) Initialize(host, user, password, dbname string) {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	var err error
	app.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter()

	app.initializeRoutes()
}

// Link all handlers to paths
func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/albums", app.getAlbums).Methods("GET")
	app.Router.HandleFunc("/album", app.createAlbum).Methods("POST")
	app.Router.HandleFunc("/album/{id:[0-9]+}", app.getAlbum).Methods("GET")
	app.Router.HandleFunc("/album/{id:[0-9]+}", app.updateAlbum).Methods("PUT")
	app.Router.HandleFunc("/album/{id:[0-9]+}", app.deleteAlbum).Methods("DELETE")
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

// Write JSON with Error to w ResponseWriter
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// Write JSON to w ResponseWriter
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Handler for getting Albums
func (app *App) getAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid album ID")
		return
	}

	a := album{ID: id}
	if err := a.getAlbum(app.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Album not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}

	respondWithJSON(w, http.StatusOK, a)
}

// Handler for getting a single Album by ID
func (app *App) getAlbums(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	albums, err := getAlbums(app.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, albums)
}

// Handler for creating an Album
func (app *App) createAlbum(w http.ResponseWriter, r *http.Request) {
	var a album
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&a); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := a.createAlbum(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, a)
}

// Handler for updating an Album by ID
func (app *App) updateAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid album ID")
		return
	}

	var a album
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&a); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	a.ID = id

	if err := a.updateAlbum(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, a)
}

// Handler for deleting an Album by ID
func (app *App) deleteAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload request")
		return
	}

	a := album{ID: id}
	if err := a.deleteAlbum(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
