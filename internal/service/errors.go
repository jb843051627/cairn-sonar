package service

import "errors"

var (
	ErrSurveyNotFound     = errors.New("survey not found")
	ErrChamberNotFound    = errors.New("chamber not found")
	ErrInstrumentNotFound = errors.New("instrument not found")
	ErrAnomalyNotFound    = errors.New("anomaly not found")
	ErrNoEcho             = errors.New("survey has no usable echo")
	ErrInvalidSurvey      = errors.New("invalid survey")
	ErrInvalidPulse       = errors.New("invalid pulse")
	ErrInvalidCalibration = errors.New("invalid calibration")
	ErrArchiveBlocked     = errors.New("survey is not ready for archive")
	ErrArchiveFailed      = errors.New("archive writer failed")
	ErrArchiveNotFound    = errors.New("archive not found")
)
