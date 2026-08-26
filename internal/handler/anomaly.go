package handler

import (
	"cairn-sonar/internal/api"
	"cairn-sonar/internal/model"
	"net/http"
	"strings"
)

func (s *Server) anomaly(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/anomaly/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		api.Error(w, 400, "bad_path", "anomaly id required")
		return
	}
	if r.Method == http.MethodPost && len(parts) > 1 && parts[1] == "review" {
		var in api.ReviewRequest
		if err := api.DecodeJSON(r, &in); err != nil {
			api.Error(w, 400, "bad_request", err.Error())
			return
		}
		if err := s.svc.ReviewAnomaly(r.Context(), parts[0], model.AnomalyState(in.Decision), in.Reviewer, in.Comment); err != nil {
			api.Error(w, 422, "review_failed", err.Error())
			return
		}
		api.JSON(w, 200, map[string]string{"status": "reviewed"})
		return
	}
	item, err := s.svc.ListAnomalies(r.Context(), parts[0])
	if err != nil {
		api.Error(w, 500, "list_failed", err.Error())
		return
	}
	api.JSON(w, 200, item)
}
