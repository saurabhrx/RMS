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
	srv.HandleFunc("/restaurants", handler.GetAllRestaurant).Methods("GET")
	srv.HandleFunc("/restaurant/{restaurant_id}/menu", handler.GetDishesByRestaurant).Methods("GET")
	srv.HandleFunc("/user/{user_id}/restaurant/{restaurant_id}/distance", handler.CalculateDistance).Methods("GET")
	protected := srv.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/logout", handler.Logout).Methods("POST")

	roleProtected := protected.NewRoute().Subrouter()
	roleProtected.Use(middleware.AuthRole("admin", "sub-admin"))
	roleProtected.HandleFunc("/create-restaurant", handler.CreateRestaurant).Methods("POST")
	roleProtected.HandleFunc("/create-dish", handler.CreateDish).Methods("POST")
	return srv
}
