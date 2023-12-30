package router

import (
	"main/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/ip-address/test", controllers.Statistics).Methods("POST")
	router.HandleFunc("/api/v1/mod_jobs/jobs", controllers.GetAllJobs).Methods("GET")
	router.HandleFunc("/api/v1/mod_jobs/test/{id}", controllers.Do_get_job).Methods("GET")
	router.HandleFunc("/api/v1/mod_jobs/job", controllers.Do_create_job).Methods("POST")
	router.HandleFunc("/api/v1/ip-mon/test", controllers.Handle_test_ip).Methods("POST")

	return router

}
