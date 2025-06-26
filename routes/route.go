package routes

import (
	"RMS/handler"
	"RMS/middleware"
	"RMS/models"
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

	// admin/subadmin

	roleProtected := protected.NewRoute().Subrouter()
	roleProtected.Use(middleware.AuthRole(models.RoleAdmin, models.RoleSubadmin))
	roleProtected.HandleFunc("/create-restaurant", handler.CreateRestaurant).Methods("POST")
	roleProtected.HandleFunc("/create-dish", handler.CreateDish).Methods("POST")
	roleProtected.HandleFunc("/get-users", handler.GetUsers).Methods("GET")
	roleProtected.HandleFunc("/user/restaurants", handler.GetRestaurantByUerID).Methods("GET")
	roleProtected.HandleFunc("/user/dish", handler.GetDishesByUserID).Methods("GET")
	roleProtected.HandleFunc("/create-user", handler.CreateUser).Methods("POST")

	// only admin

	adminOnly := protected.NewRoute().Subrouter()
	adminOnly.Use(middleware.AuthRole(models.RoleAdmin))
	adminOnly.HandleFunc("/get-subadmin", handler.GetAllSubadmin).Methods("GET")

	return srv
}
