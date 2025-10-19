package pkg

import "sync"

type MessageQueue struct {
	messages []string
	capacity int
	start    int
	count    int
	mu       sync.Mutex
}

func NewMessageQueue(capacity int) *MessageQueue {
	return &MessageQueue{
		messages: make([]string, capacity),
		capacity: capacity,
	}
}

func (q *MessageQueue) Add(msg string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	index := (q.start + q.count) % q.capacity
	q.messages[index] = msg

	if q.count < q.capacity {
		q.count++
	} else {
		// overwrite oldest
		q.start = (q.start + 1) % q.capacity
	}
}

func (q *MessageQueue) GetAll() []string {
	q.mu.Lock()
	defer q.mu.Unlock()

	result := make([]string, q.count)
	for i := 0; i < q.count; i++ {
		index := (q.start + i) % q.capacity
		result[i] = q.messages[index]
	}
	return result
}
