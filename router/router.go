package router

import (
	"github.com/gorilla/mux"
	"github.com/mateuscdiniz/simplecrud/middleware"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/job/{id}", middleware.GetJob).Methods("GET")
	router.HandleFunc("/api/job", middleware.GetAllJobs).Methods("GET")
	router.HandleFunc("/api/job", middleware.CreateJob).Methods("POST")
	router.HandleFunc("/api/job/{id}", middleware.UpdateJob).Methods("PUT")
	router.HandleFunc("/api/job/{id}", middleware.DeleteJob).Methods("DELETE")

	return router
}
