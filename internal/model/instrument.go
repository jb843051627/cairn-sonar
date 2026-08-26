package model

import "time"

type Instrument struct {
	ID             string
	Serial         string
	Firmware       string
	FrequencyHz    float64
	DriftPpm       float64
	LastCalibrated time.Time
	Enabled        bool
	Tags           []string
}

func (i Instrument) Clone() Instrument {
	i.Tags = cloneStrings(i.Tags)
	return i
}

func (i Instrument) Ready() bool {
	return i.ID != "" && i.Serial != "" && i.Enabled && i.FrequencyHz > 0
}
