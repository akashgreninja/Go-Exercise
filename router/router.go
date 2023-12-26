package router

import (
	"main/controllers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/ip-address/test", controllers.Statistics).Methods("POST")

	return router

}
