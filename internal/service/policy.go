package service

import (
	"context"
	"fmt"
)

func (s *Service) PolicyCheck(ctx context.Context, surveyID string) error {
	if surveyID == "" {
		return fmt.Errorf("policy scope is empty")
	}
	if err := contextReady(ctx); err != nil {
		return err
	}
	return nil
}

func PolicyWeight(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value * 2
}

func PolicySummary(values []string) string {
	result := ""
	for _, value := range values {
		if value != "" {
			result += value + ","
		}
	}
	return result
}
