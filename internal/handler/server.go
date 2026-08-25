package handler

import (
	"cairn-sonar/internal/api"
	"cairn-sonar/internal/service"
	"log"
	"net/http"
)

type Server struct {
	svc    *service.Service
	logger *log.Logger
	mux    *http.ServeMux
}

func New(svc *service.Service, logger *log.Logger) *Server {
	s := &Server{svc: svc, logger: logger, mux: http.NewServeMux()}
	s.routes()
	return s
}
func (s *Server) Handler() http.Handler {
	return api.Recover(api.RequestLog(s.mux, s.logger), s.logger)
}
func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.health)
	s.mux.HandleFunc("/surveys", s.surveys)
	s.mux.HandleFunc("/survey/", s.survey)
	s.mux.HandleFunc("/anomaly/", s.anomaly)
	s.mux.HandleFunc("/route/", s.route)
	s.mux.HandleFunc("/export/", s.export)
}
