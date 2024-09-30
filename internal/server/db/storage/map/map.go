package _map

import (
	"fmt"
	"sync"
)

type MapStorage struct {
	mtx sync.Mutex
	m   map[string]string
}

func (s *MapStorage) New() {
	s.m = make(map[string]string)
}

func (s *MapStorage) Get(key string) (string, error) {
	s.mtx.Lock()
	val, ok := s.m[key]
	s.mtx.Unlock()
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}

	return val, nil
}

func (s *MapStorage) Set(key, value string) error {
	s.mtx.Lock()
	s.m[key] = value
	s.mtx.Unlock()
	return nil
}

func (s *MapStorage) Del(key string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.m[key]
	if !ok {
		return fmt.Errorf("key %s not found", key)
	}

	delete(s.m, key)

	return nil
}

func (s *MapStorage) Close() error { return nil }
