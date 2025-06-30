package routes

import (
	"RMS/handler"
	"RMS/middleware"
	"RMS/models"
	"github.com/gorilla/mux"
)

func SetupTodoRoutes() *mux.Router {
	srv := mux.NewRouter()
	api := srv.PathPrefix("/api/v1").Subrouter()
	// public route
	api.HandleFunc("/register", handler.Register).Methods("POST")
	api.HandleFunc("/login", handler.Login).Methods("POST")
	api.HandleFunc("/refresh", handler.Refresh).Methods("POST")
	api.HandleFunc("/restaurants", handler.GetAllRestaurant).Methods("GET")
	api.HandleFunc("/restaurant/{restaurant_id}/menu", handler.GetDishesByRestaurant).Methods("GET")
	api.HandleFunc("/user/{address_id}/restaurant/{restaurant_id}/distance", handler.CalculateDistance).Methods("GET")

	protected := api.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// private route
	protected.HandleFunc("/create-address", handler.CreateAddress).Methods("POST")
	protected.HandleFunc("/logout", handler.Logout).Methods("POST")

	// admin/sub-admin
	roleProtected := protected.NewRoute().Subrouter()
	roleProtected.Use(middleware.AuthRole(models.RoleAdmin, models.RoleSubadmin))
	roleProtected.HandleFunc("/create-restaurant", handler.CreateRestaurant).Methods("POST")
	roleProtected.HandleFunc("/create-dish", handler.CreateDish).Methods("POST")
	roleProtected.HandleFunc("/user/restaurants", handler.GetRestaurantByUerID).Methods("GET")
	roleProtected.HandleFunc("/user/dish", handler.GetDishesByUserID).Methods("GET")
	roleProtected.HandleFunc("/get-users", handler.GetUsers).Methods("GET")

	// only admin
	adminOnly := protected.NewRoute().Subrouter()
	adminOnly.Use(middleware.AuthRole(models.RoleAdmin))
	adminOnly.HandleFunc("/admin/get-subadmin", handler.GetAllSubadmin).Methods("GET")
	adminOnly.HandleFunc("/admin/create-user", handler.CreateUserByAdmin).Methods("POST")

	//only sub-admin
	subadminOnly := protected.NewRoute().Subrouter()
	subadminOnly.Use(middleware.AuthRole(models.RoleSubadmin))
	subadminOnly.HandleFunc("/subadmin/create-user", handler.CreateUserBySubadmin).Methods("POST")

	return srv
}
