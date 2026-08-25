package handler

import (
	"cairn-sonar/internal/api"
	"net/http"
	"strings"
)

func (s *Server) export(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/export/")
	data, err := s.svc.ExportSurvey(r.Context(), id)
	if err != nil {
		api.Error(w, 422, "export_failed", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
