package storage

import (
	"custom-in-memory-db/internal/server/cmd"
	"log/slog"
)

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Del(key string) error
	Recover(conf cmd.Config, lg *slog.Logger) error
	Close() error
}
