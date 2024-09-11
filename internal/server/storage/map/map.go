package _map

import (
	"custom-in-memory-db/internal/server/parser"
	"fmt"
)

type MapStorage struct {
	m map[string]string
}

func (s *MapStorage) New() {
	s.m = make(map[string]string)
}

// At this point only valid parser.Command struct present

func (s *MapStorage) Get(c parser.Command) (string, error) {
	val, ok := s.m[c.Args[0]]
	if !ok {
		return "", fmt.Errorf("key %s not found", c.Args[0])
	}

	return val, nil
}

func (s *MapStorage) Set(c parser.Command) error {
	s.m[c.Args[0]] = c.Args[1]

	return nil
}

func (s *MapStorage) Del(c parser.Command) error {
	_, ok := s.m[c.Args[0]]
	if !ok {
		return fmt.Errorf("key %s not found", c.Args[0])
	}

	delete(s.m, c.Args[0])

	return nil
}
