package handler

import (
	"cairn-sonar/internal/api"
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/service"
	"errors"
	"net/http"
	"strings"
)

func (s *Server) surveys(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var in api.SurveyRequest
		if err := api.DecodeJSON(r, &in); err != nil {
			api.Error(w, 400, "bad_request", err.Error())
			return
		}
		err := s.svc.CreateSurvey(r.Context(), model.Survey{ID: in.ID, ChamberID: in.ChamberID, Lead: in.Lead, Notes: in.Notes})
		if err != nil {
			api.Error(w, 422, "create_failed", err.Error())
			return
		}
		api.JSON(w, 201, map[string]string{"id": in.ID})
		return
	}
	items, err := s.svc.ListSurveys(r.Context(), model.SurveyStatus(r.URL.Query().Get("status")))
	if err != nil {
		api.Error(w, 500, "list_failed", err.Error())
		return
	}
	api.JSON(w, 200, items)
}
func (s *Server) survey(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/survey/")
	if r.Method == http.MethodPost && strings.HasSuffix(id, "/start") {
		id = strings.TrimSuffix(id, "/start")
		if err := s.svc.StartSurvey(r.Context(), id); err != nil {
			api.Error(w, 422, "start_failed", err.Error())
			return
		}
		api.JSON(w, 200, map[string]string{"status": "active"})
		return
	}
	item, err := s.svc.GetSurvey(r.Context(), id)
	if errors.Is(err, service.ErrSurveyNotFound) {
		api.Error(w, 404, "not_found", err.Error())
		return
	}
	if err != nil {
		api.Error(w, 500, "get_failed", err.Error())
		return
	}
	api.JSON(w, 200, item)
}
