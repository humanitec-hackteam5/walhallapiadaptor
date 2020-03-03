package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"humanitec.io/walhallapiadaptor/internal/walhallapi"
)

type server struct {
	router       http.Handler
	newWalhall   func(jwt string) (walhallapi.WalhallAPIer, error)
	registryName string
}

func main() {
	var s server

	walhallAPIPrefix := os.Getenv("WALHALL_API_PREFIX")
	var reusableClient http.Client

	s.newWalhall = func(jwt string) (walhallapi.WalhallAPIer, error) {
		return walhallapi.New(walhallAPIPrefix, jwt, &reusableClient)
	}

	s.registryName = os.Getenv("WALHALL_REGISTRY")

	log.Println("Setting up Routes")
	s.setupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on Port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, s.router)))
}
