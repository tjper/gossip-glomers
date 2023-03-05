package safe

import "sync"

func NewSlice[T comparable]() *Slice[T] {
	return &Slice[T]{
		mutex: new(sync.RWMutex),
		items: make([]T, 0),
	}
}

type Slice[T comparable] struct {
	mutex *sync.RWMutex
	items []T
}

func (s *Slice[T]) Contains(item T) bool {
	s.mutex.RLock()
  defer s.mutex.RUnlock()

	for i := 0; i < len(s.items); i++ {
		if s.items[i] == item {
			return true
		}
	}
	return false
}

func (s *Slice[T]) Length() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.items)
}

func (s *Slice[T]) Push(items ...T) {
	s.mutex.Lock()
	s.items = append(s.items, items...)
	s.mutex.Unlock()
}

func (s *Slice[T]) Lpop() T {
	s.mutex.Lock()
	item := s.items[0]
	s.items = s.items[1:]
	s.mutex.Unlock()

	return item
}

func (s *Slice[T]) Set(newItems ...T) {
	items := make([]T, len(newItems))
	copy(items, newItems)

	s.mutex.Lock()
	s.items = items
	s.mutex.Unlock()
}

func (s *Slice[T]) Contents() []T {
	s.mutex.RLock()
	items := make([]T, len(s.items))
	copy(items, s.items)
	s.mutex.RUnlock()
	return items
}
