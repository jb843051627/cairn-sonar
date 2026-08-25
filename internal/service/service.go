package service

import (
	"cairn-sonar/internal/clock"
	"cairn-sonar/internal/model"
	"cairn-sonar/internal/rules"
	"cairn-sonar/internal/store"
	"cairn-sonar/internal/worker"
	"context"
	"fmt"
	"sync"
	"time"
)

type ArchiveWriter interface {
	Write(context.Context, model.Survey, []model.EchoProfile) (model.Archive, error)
}

type Service struct {
	repo         *store.Repository
	now          clock.Clock
	thresholds   rules.Thresholds
	queue        *worker.Queue
	archive      ArchiveWriter
	cacheMu      sync.RWMutex
	echoCache    map[string][]model.EchoProfile
	anomalyCache map[string][]model.Anomaly
}

type Config struct {
	Clock     clock.Clock
	Archive   ArchiveWriter
	QueueSize int
}

func New(repo *store.Repository, cfg Config) *Service {
	c := cfg.Clock
	if c == nil {
		c = clock.Real{}
	}
	q := worker.NewQueue(cfg.QueueSize)
	return &Service{repo: repo, now: c, thresholds: rules.DefaultThresholds(), queue: q, archive: cfg.Archive, echoCache: make(map[string][]model.EchoProfile), anomalyCache: make(map[string][]model.Anomaly)}
}

func (s *Service) Repository() *store.Repository { return s.repo }
func (s *Service) PendingWork() int              { return s.queue.Pending() }
func (s *Service) Close()                        { s.queue.Close() }

func (s *Service) mapStoreError(err error, target error) error {
	if err == nil {
		return nil
	}
	if err == store.ErrNotFound {
		return target
	}
	return fmt.Errorf("service storage: %w", err)
}

func contextReady(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
func utcNow(c clock.Clock) time.Time { return c.Now().UTC() }
