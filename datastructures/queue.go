package datastructures

import (
	"errors"
	"net/url"
	"sync"
)

var (
	// ErrEmptyQueue is returned when there's no elements left in the queue
	ErrEmptyQueue = errors.New("Queue is empty")
)

type Queue struct {
	arr []*url.URL
	mu  sync.Mutex
}

// NewQueue returns a new queue object (which is based on a slice working as a list)
func NewQueue() *Queue {
	return &Queue{
		arr: make([]*url.URL, 0),
	}
}

// Dequeue removes a URL reference from the queue and it returns it back to the caller.
// It errors if you try to Dequeue() on an empty queue.
func (q *Queue) Dequeue() (*url.URL, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.arr) == 0 {
		return nil, ErrEmptyQueue
	}

	el := q.arr[0]

	// Delete q.arr[0] from the list. Doing it like this so that the internal
	// slice structure is able to stop referencing the value and the GC can then
	// collect it.
	copy(q.arr[0:], q.arr[1:])
	q.arr[len(q.arr)-1] = nil
	q.arr = q.arr[:len(q.arr)-1]

	return el, nil
}

// Enqueue puts a URL reference into the queue
func (q *Queue) Enqueue(el *url.URL) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.arr = append(q.arr, el)
}
