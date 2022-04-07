package app

import (
	"bus-router-backend/mvc/controllers"
	"bus-router-backend/mvc/middleware"
	"bus-router-backend/mvc/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartApp() {
	db := services.MongoClient
	fmt.Println(db)
	r := mux.NewRouter()
	r.HandleFunc("/admins", controllers.GetAdmin)
	r.HandleFunc("/adminLogin", controllers.AdminLogin).Methods("POST")
	r.Handle("/get/admin", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.AdminInside)))
	r.Handle("/logout/admin", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.AdminLogout))).Methods("POST")
	r.HandleFunc("/signUp", controllers.SignUp).Methods("POST")
	r.HandleFunc("/signIn", controllers.SignIn).Methods("POST")
	r.Handle("/inside/user", middleware.MiddlewareWrapper(middleware.IsAuth(controllers.GetInside))).Methods("GET")
	r.Handle("/logout/user", middleware.MiddlewareWrapper(middleware.IsAuth(controllers.LogoutForUser))).Methods("POST")
	r.Handle("/get/users", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.GetAllUsers)))
	r.Handle("/get/stations", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.GetAllStations)))
	r.Handle("/delete/station", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.DeleteStation))).Methods("DELETE")
	r.Handle("/delete/voyage", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.DeleteVoyage))).Methods("DELETE")
	r.Handle("/add/station", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.AddStation))).Methods("POST")
	r.Handle("/update/station", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.UpdatePassengerOfStation))).Methods("POST")
	r.Handle("/add/voyage", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.AddVoyage))).Methods("POST")
	r.Handle("/get/voyages", middleware.MiddlewareForAdmin(middleware.IsAdminLoggedIn(controllers.GetAllVoyages)))
	log.Fatal(http.ListenAndServe(":8080", r))
}
