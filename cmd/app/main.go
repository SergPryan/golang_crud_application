package main

import (
	"crud_application/internal/controller"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /vacancy/{id}", controller.HandlerGetVacancy)
	router.HandleFunc("POST /request", controller.HandlerPostRequestNewVacancies)

	log.Fatal(http.ListenAndServe(":8080", router))
}
