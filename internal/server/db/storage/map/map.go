package _map

import (
	"fmt"
	"sync"
)

type Storage struct {
	mtx sync.Mutex
	m   map[string]string
}

// New used to initialize Storage.
// Any initializations after the first one won't take effect
func New() *Storage {
	st := Storage{}
	st.m = make(map[string]string)
	return &st
}

func (s *Storage) Get(key string) (string, error) {
	s.mtx.Lock()
	val, ok := s.m[key]
	s.mtx.Unlock()
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}

	return val, nil
}

func (s *Storage) Set(key, value string) error {
	s.mtx.Lock()
	s.m[key] = value
	s.mtx.Unlock()

	return nil
}

func (s *Storage) Del(key string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.m[key]
	if !ok {
		return fmt.Errorf("key %s not found", key)
	}

	delete(s.m, key)

	return nil
}
