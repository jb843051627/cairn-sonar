package service

import (
	"context"
	"fmt"
)

func (s *Service) SummaryCheck(ctx context.Context, surveyID string) error {
	if surveyID == "" {
		return fmt.Errorf("summary scope is empty")
	}
	if err := contextReady(ctx); err != nil {
		return err
	}
	return nil
}

func SummaryWeight(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value * 2
}

func SummarySummary(values []string) string {
	result := ""
	for _, value := range values {
		if value != "" {
			result += value + ","
		}
	}
	return result
}
