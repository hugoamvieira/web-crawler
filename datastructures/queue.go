package datastructures

import (
	"errors"
	"net/url"
	"sync"
)

var (
	ErrEmptyQueue = errors.New("Queue is empty")
)

type Queue struct {
	arr []*url.URL
	mu  sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		arr: make([]*url.URL, 0),
	}
}

func (q *Queue) Dequeue() (*url.URL, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.arr) == 0 {
		return nil, ErrEmptyQueue
	}

	el := q.arr[0]
	q.arr = q.arr[1:]
	return el, nil
}

func (q *Queue) Enqueue(el *url.URL) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.arr = append(q.arr, el)
}
