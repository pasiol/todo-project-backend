package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshalling payload failed: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	bytes, err := w.Write(response)
	if err != nil {
		log.Printf("writing response failed: %s", err)
	}
	log.Printf("response bytes %d, %v", bytes, payload)
}

func (a *App) getTodos(w http.ResponseWriter, _ *http.Request) {
	todos, err := a.searchTodos()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, todos)
}

func (a *App) getHealth(w http.ResponseWriter, _ *http.Request) {
	err := a.Pool.Ping()
	if err == nil {
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "ok"})
		return
	} else {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *App) postTodo(w http.ResponseWriter, r *http.Request) {
	var t Todo

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed todo object")
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("closing request failed: %s", err)
		}
	}(r.Body)

	if err := a.insertTodo(t); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, "new todo task created")
}

func (a *App) putTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf(" %v", vars["id"])
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	err = a.updateTodoDone(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, id)
}
