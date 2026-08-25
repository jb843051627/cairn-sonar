package model

import "time"

type Archive struct {
	ID          string
	SurveyID    string
	ObjectKey   string
	Digest      string
	SizeBytes   int64
	CompletedAt time.Time
	Verified    bool
}

func (a Archive) Valid() bool {
	return a.ID != "" && a.SurveyID != "" && a.ObjectKey != "" && a.Digest != ""
}
