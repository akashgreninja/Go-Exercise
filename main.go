package main

import (
	"fmt"
	"log"
	_ "main/database"
	"main/router"
	"net/http"
)

func main() {

	fmt.Println("started server ")
	r := router.Router()

	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening at port 4000")

	//cannot execute ping directly in golang runtime using exec

}
