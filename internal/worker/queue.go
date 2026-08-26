package worker

import "sync"

type Queue struct {
	mu      sync.RWMutex
	pending map[string]int
	jobs    chan string
	stop    chan struct{}
	once    sync.Once
}

func NewQueue(size int) *Queue {
	if size < 1 {
		size = 16
	}
	q := &Queue{pending: make(map[string]int), jobs: make(chan string, size), stop: make(chan struct{})}
	go q.loop()
	return q
}
func (q *Queue) Submit(key string) {
	q.mu.Lock()
	q.pending[key]++
	q.mu.Unlock()
	select {
	case q.jobs <- key:
	default:
	}
}
func (q *Queue) Pending() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	total := 0
	for _, n := range q.pending {
		total += n
	}
	return total
}
func (q *Queue) Complete(key string) {
	q.mu.Lock()
	if q.pending[key] > 1 {
		q.pending[key]--
	} else {
		delete(q.pending, key)
	}
	q.mu.Unlock()
}
func (q *Queue) Close() { q.once.Do(func() { close(q.stop) }) }
func (q *Queue) loop() {
	for {
		select {
		case key := <-q.jobs:
			q.Complete(key)
		case <-q.stop:
			return
		}
	}
}
