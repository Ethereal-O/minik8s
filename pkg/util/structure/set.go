package structure

import "sync"

// Set : A thread-safe set
type Set struct {
	lock   sync.Mutex
	record map[any]struct{}
}

func (s *Set) Init() {
	s.record = make(map[any]struct{})
}

func (s *Set) Put(element any) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	// element is in s.record
	if _, ok := s.record[element]; ok {
		return false
	} else {
		s.record[element] = struct{}{}
		return true
	}
}

func (s *Set) Del(element any) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	// element is in s.record
	if _, ok := s.record[element]; ok {
		delete(s.record, element)
		return true
	} else {
		return false
	}
}

func (s *Set) Contain(element any) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.record[element]
	return ok
}

func (s *Set) GetAll() []any {
	var allElements []any
	for element := range s.record {
		allElements = append(allElements, element)
	}
	return allElements
}
