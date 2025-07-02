package routes

import (
	"RMS/handler"
	"RMS/middleware"
	"RMS/models"
	"RMS/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupTodoRoutes() *mux.Router {
	srv := mux.NewRouter()

	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseError(w, http.StatusOK, "server is running")
	})

	api := srv.PathPrefix("/api/v1").Subrouter()
	// public route
	api.HandleFunc("/user/register", handler.RegisterUser).Methods("POST")
	api.HandleFunc("/user/login", handler.LoginUser).Methods("POST")
	api.HandleFunc("/admin/login", handler.LoginAdmin).Methods("POST")
	api.HandleFunc("/subadmin/login", handler.LoginSubadmin).Methods("POST")
	api.HandleFunc("/user/login", handler.LoginUser).Methods("POST")
	api.HandleFunc("/refresh", handler.Refresh).Methods("POST")
	api.HandleFunc("/restaurant", handler.GetAllRestaurant).Methods("GET")
	api.HandleFunc("/restaurant/{restaurant_id}/menu", handler.GetDishesByRestaurant).Methods("GET")

	protected := api.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// private route
	protected.HandleFunc("/address", handler.CreateAddress).Methods("POST")
	//protected.HandleFunc("/logout", handler.Logout).Methods("POST")
	protected.HandleFunc("/user/{address_id}/restaurant/{restaurant_id}/distance", handler.CalculateDistance).Methods("GET")

	// admin/subadmin
	roleProtected := protected.NewRoute().Subrouter()
	roleProtected.Use(middleware.AuthRole(models.RoleAdmin, models.RoleSubadmin))
	roleProtected.HandleFunc("/restaurant", handler.CreateRestaurant).Methods("POST")
	roleProtected.HandleFunc("/dish", handler.CreateDish).Methods("POST")
	roleProtected.HandleFunc("/user/restaurant", handler.GetRestaurantByUserID).Methods("GET")
	roleProtected.HandleFunc("/user/dish", handler.GetDishesByUserID).Methods("GET")
	roleProtected.HandleFunc("/users", handler.GetUsers).Methods("GET")

	// only admin
	adminOnly := protected.NewRoute().Subrouter()
	adminOnly.Use(middleware.AuthRole(models.RoleAdmin))
	adminOnly.HandleFunc("/admin/subadmin", handler.GetAllSubadmin).Methods("GET")
	adminOnly.HandleFunc("/admin/user", handler.CreateUserByAdmin).Methods("POST")

	//only subadmin
	subadminOnly := protected.NewRoute().Subrouter()
	subadminOnly.Use(middleware.AuthRole(models.RoleSubadmin))
	subadminOnly.HandleFunc("/subadmin/user", handler.CreateUserBySubadmin).Methods("POST")

	return srv
}
