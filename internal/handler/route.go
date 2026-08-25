package handler

import (
	"cairn-sonar/internal/api"
	"net/http"
	"strings"
)

func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/route/")
	if r.Method == http.MethodPost {
		route, err := s.svc.PlanRoute(r.Context(), id)
		if err != nil {
			api.Error(w, 422, "route_failed", err.Error())
			return
		}
		api.JSON(w, 201, route)
		return
	}
	route, err := s.svc.GetRoute(r.Context(), id)
	if err != nil {
		api.Error(w, 404, "not_found", err.Error())
		return
	}
	api.JSON(w, 200, route)
}
