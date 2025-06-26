package routes

import (
	"RMS/handler"
	"RMS/middleware"
	"github.com/gorilla/mux"
)

func SetupTodoRoutes() *mux.Router {
	srv := mux.NewRouter()
	srv.HandleFunc("/register", handler.Register).Methods("POST")
	srv.HandleFunc("/login", handler.Login).Methods("POST")

	protected := srv.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/logout", handler.Logout).Methods("POST")
	return srv
}
