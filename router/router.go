package router

import "github.com/gorilla/mux"

func Router() *mux.Router {
	router := mux.NewRouter()
	// router.HandleFunc("/api/v1/ip-address/test")

	return router

}
