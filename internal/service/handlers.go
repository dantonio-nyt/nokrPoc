package service

import (
	"net/http"
)

func (h *HermesService) HealthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("system is up"))
	})
}
