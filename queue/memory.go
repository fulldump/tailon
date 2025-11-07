package queue

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

type MemoryService struct {
	Queues      map[string]Queue // todo: replace by sync.Map
	QueuesMutex sync.RWMutex
}

func NewMemoryService() *MemoryService {
	return &MemoryService{
		Queues: map[string]Queue{},
	}
}

func (m *MemoryService) GetQueue(name string) (Queue, error) {

	m.QueuesMutex.RLock()
	defer m.QueuesMutex.RUnlock()

	q, exists := m.Queues[name]
	if !exists {
		return nil, fmt.Errorf("queue '%s' does not exist", name)
	}

	return q, nil
}

func (m *MemoryService) ListQueues() ([]string, error) {

	m.QueuesMutex.RLock()
	defer m.QueuesMutex.RUnlock()

	result := make([]string, 0, len(m.Queues))

	for name := range m.Queues {
		result = append(result, name)
	}

	sort.Strings(result)

	return result, nil
}

func (m *MemoryService) CreateQueue(name string) (Queue, error) {

	m.QueuesMutex.Lock()
	defer m.QueuesMutex.Unlock()

	if _, exists := m.Queues[name]; exists {
		return nil, fmt.Errorf("queue '%s' already exists", name)
	}

	q := NewMemoryQueue()
	m.Queues[name] = q

	return q, nil
}

func (m *MemoryService) DeleteQueue(name string) error {

	m.QueuesMutex.Lock()
	defer m.QueuesMutex.Unlock()

	if _, exists := m.Queues[name]; !exists {
		return fmt.Errorf("queue '%s' does not exist", name)
	}

	delete(m.Queues, name)
	return nil
}

type MemoryQueue struct {
	Queue  chan JSON
	Writes int64
	Reads  int64
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		Queue: make(chan JSON, 10*1000*1000),
	}
}

func (m *MemoryQueue) Write(item JSON) error {

	atomic.AddInt64(&m.Writes, 1)

	m.Queue <- item // todo: make it sync?

	return nil
}

func (m *MemoryQueue) Read() (JSON, error) {

	item := <-m.Queue

	atomic.AddInt64(&m.Reads, 1)

	return item, nil
}
