package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/nokrPOC/internal/config"
	"github.com/nokrPOC/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	serviceConfig, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("couldn't initialize the config: %s", err))
	}

	HermesService, err := service.NewService(context.Background(), serviceConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to start service %v", err))
	}

	router := mux.NewRouter()
	HermesService.SetupRoutes(router)
	log.Printf("Listing on port: %s", serviceConfig.Port)
	err = http.ListenAndServe(":"+ serviceConfig.Port, router)
	if err != nil {
		panic(fmt.Sprintf("unable to serve %s", err))
	}
}
