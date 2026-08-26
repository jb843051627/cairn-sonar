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
// Close signals the consumer loop to exit without closing the jobs channel.
// The jobs channel is the sender's side of the contract: Submit may still be
// racing with shutdown (an in-flight collection request can land in the queue
// after Close is called). Closing jobs here would turn that concurrent send
// into a "send on closed channel" panic. The loop drains via the stop signal
// alone, so jobs must remain open and rely on GC once the queue is abandoned.
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
