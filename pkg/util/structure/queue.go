package structure

import "sync"

// Queue : A thread-safe queue
type Queue struct {
	items []any
	lock  sync.Mutex
}

func (q *Queue) Init() {
	q.items = make([]any, 0)
}

func (q *Queue) Push(item any) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.items = append(q.items, item)
}

func (q *Queue) Pop() any {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) == 0 {
		return nil
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue) Front() any {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) == 0 {
		return nil
	}

	return q.items[0]
}
