package service

import "github.com/gorilla/mux"

func (h *HermesService) SetupRoutes(router *mux.Router) {
	router.Handle("/healthcheck", h.HealthCheck()).Methods("GET")
}