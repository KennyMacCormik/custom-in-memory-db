package storage

import (
	"custom-in-memory-db/internal/server/parser"
)

type Storage interface {
	Get(c parser.Command) (string, error)
	Set(c parser.Command) error
	Del(c parser.Command) error
}
