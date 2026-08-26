package api

import "time"

type SurveyRequest struct {
	ID        string   `json:"id"`
	ChamberID string   `json:"chamber_id"`
	Lead      string   `json:"lead"`
	Notes     []string `json:"notes"`
}

type PulseRequest struct {
	ID           string    `json:"id"`
	SurveyID     string    `json:"survey_id"`
	InstrumentID string    `json:"instrument_id"`
	Sequence     int       `json:"sequence"`
	EmittedAt    time.Time `json:"emitted_at"`
	DurationMS   int       `json:"duration_ms"`
	GainDB       float64   `json:"gain_db"`
	Samples      []float64 `json:"samples"`
}

type ReviewRequest struct {
	Decision string `json:"decision"`
	Comment  string `json:"comment"`
	Reviewer string `json:"reviewer"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}
