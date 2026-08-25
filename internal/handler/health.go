package handler

import (
	"cairn-sonar/internal/api"
	"net/http"
	"time"
)

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	api.JSON(w, http.StatusOK, api.HealthResponse{Status: "ok", Time: time.Now().UTC().Format(time.RFC3339Nano)})
}
