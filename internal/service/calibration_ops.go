package service

import (
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/store"
	"context"
	"errors"
	"fmt"
	"math"
)

func (s *Service) CalibrateInstrument(ctx context.Context, c model.Calibration) error {
	if err := contextReady(ctx); err != nil {
		return err
	}
	if !c.Valid() {
		return ErrInvalidCalibration
	}
	instrument, err := s.repo.GetInstrument(ctx, c.InstrumentID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return ErrInstrumentNotFound
		}
		return fmt.Errorf("load instrument: %w", err)
	}
	c.OffsetDB = c.ReferenceDB - c.MeasuredDB
	c.Passed = math.Abs(c.OffsetDB) <= 2
	if err := s.repo.SaveCalibration(ctx, c); err != nil {
		return err
	}
	instrument.DriftPpm += c.OffsetDB
	instrument.LastCalibrated = utcNow(s.now)
	if !c.Passed {
		instrument.Enabled = false
	}
	return s.repo.PutInstrument(ctx, instrument)
}
func (s *Service) CalibrationCount(ctx context.Context, id string) (int, error) {
	return s.repo.CountCalibrations(ctx, id)
}
