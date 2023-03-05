package safe

import "sync"

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		mutex: new(sync.RWMutex),
		hash:  make(map[T]struct{}),
	}
}

type Set[T comparable] struct {
	mutex *sync.RWMutex
	hash  map[T]struct{}
}

func (s *Set[T]) Add(item T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.hash[item]
	if ok {
		return false
	}

	s.hash[item] = struct{}{}
	return true
}

func (s *Set[T]) Contents() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	items := make([]T, 0, len(s.hash))
	for item := range s.hash {
		items = append(items, item)
	}
	return items
}
