package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Reading environment failed.")
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/todos", a.getTodos).Methods("GET")
	a.Router.HandleFunc("/todos", a.postTodo).Methods("POST")
}

func (a *App) Run() {

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	origins := handlers.AllowedOrigins([]string{os.Getenv("ALLOWED_ORIGINS")})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodOptions, http.MethodConnect, http.MethodPost})
	maxAge := handlers.MaxAge(60)

	address := fmt.Sprintf("0.0.0.0:%s", os.Getenv("APP_PORT"))
	server := &http.Server{
		Addr:    address,
		Handler: handlers.CORS(headers, origins, methods, maxAge)(a.Router),
	}

	log.Printf("starting REST-backend in %s.", address)
	log.Printf("Version: %s , build: %s", Version, Build)
	log.Printf("Allowed origins: %s", os.Getenv("ALLOWED_ORIGINS"))
	log.Fatal(server.ListenAndServe())
}
