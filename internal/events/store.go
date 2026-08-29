package events

import "sync"

type Store struct {
	mu sync.RWMutex
	events []StreamEvent
	max int
}

func NewStore(max int) *Store {
	return &Store {
		events: make([]StreamEvent, 0, max),
		max: max,
	}
}

func (s *Store) Add(event StreamEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append([]StreamEvent{event}, s.events...)

	if len(s.events) > s.max {
		s.events = s.events[:s.max]
	}
}

func (s *Store) List() []StreamEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]StreamEvent, len(s.events))
	copy(result, s.events)

	return result
}