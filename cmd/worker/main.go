package main

import (
	"context"
	"fmt"
	"github.com/nokrPOC/internal/config"
	"github.com/nokrPOC/internal/service"
	"log"
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

	err = HermesService.StartListeningForMessages()
	if err != nil {
		panic(fmt.Sprintf("failed to start listening %v", err))
	}
}
