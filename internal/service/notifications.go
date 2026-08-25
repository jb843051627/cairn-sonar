package service

import (
	"context"
	"fmt"
)

func (s *Service) NotificationCheck(ctx context.Context, surveyID string) error {
	if surveyID == "" {
		return fmt.Errorf("notification scope is empty")
	}
	if err := contextReady(ctx); err != nil {
		return err
	}
	return nil
}

func NotificationWeight(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value * 2
}

func NotificationSummary(values []string) string {
	result := ""
	for _, value := range values {
		if value != "" {
			result += value + ","
		}
	}
	return result
}
