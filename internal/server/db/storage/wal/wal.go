package wal

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/seg"
	"custom-in-memory-db/internal/server/db/storage"
	"errors"
	atomicUber "go.uber.org/atomic"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

// ErrWalWriteFailed is the error returned by w.write when batch failed to write to wal.
// (w.write must return ErrWalWriteFailed itself, not an error wrapping ErrWalWriteFailed, because callers will test for EOF using ==.)
var ErrWalWriteFailed = errors.New("wal write failed")
var writeOk = errors.New("ok")

// Storage is the same as map.Storage and adds wal implementation.
// Expect files to be written to a WAL_SEG_PATH
type Storage struct {
	st storage.Storage

	batch     atomicUber.String
	batchMax  int32
	batchSize atomic.Int32
	// Timer implementation.
	// We're flipping a coin every flipTimer duration
	flipTimer time.Duration
	coin      atomic.Bool
	closer    chan struct{}

	writeHappens atomic.Bool
	writer       io.WriteCloser
}

// Set sets provided value for the provided key.
// Set is thread-safe
func (s *Storage) Set(key, value string) error {
	err := s.lockOrWrite("SET", key, value)
	if err != nil && err != writeOk {
		return err
	}
	return s.st.Set(key, value)
}

// Del removes provided key and it's value.
// Del is thread-safe
func (s *Storage) Del(key string) error {
	err := s.lockOrWrite("DEL", key)
	if err != nil && err != writeOk {
		return err
	}
	return s.st.Del(key)
}

// Get returns a value of the provided key.
func (s *Storage) Get(key string) (string, error) {
	return s.st.Get(key)
}

// New used to initialize Storage.
// Any initializations after the first one won't take effect
func New(conf cmd.Config, st storage.Storage, lg *slog.Logger) (*Storage, error) {
	s := Storage{}
	s.st = st
	s.batchMax = int32(conf.Wal.BatchMax)
	s.flipTimer = conf.Wal.BatchTimeout
	s.writeHappens.Store(false)
	s.closer = make(chan struct{})
	go s.coinFlipper(s.closer)
	sg, err := seg.New(conf)
	if err != nil {
		return nil, err
	}
	if conf.Wal.Recover {
		err = Recover(sg, st.Set, st.Del, lg)
	}
	s.writer = sg

	return &s, nil
}

// Close gracefully stops the Storage
func (s *Storage) Close() error {
	// stop ticker and its goroutine
	s.closer <- struct{}{}
	<-s.closer
	// wait for wal to be written
	for s.batch.Load() != "" {
		runtime.Gosched()
	}
	return s.writer.Close()
}

// addToBuff loads command to batch and increments batchSize.
func (s *Storage) addToBuff(args ...string) {
	cmnd := strings.Join(append(args[:], "\n"), " ")

	for {
		oldVal := s.batch.Load()
		if !s.writeHappens.Load() && s.batchSize.Load() < s.batchMax && s.batch.CompareAndSwap(oldVal, strings.Join([]string{oldVal, cmnd}, "")) {
			// This doesn't ensure strict batch size. Overflow might happen
			s.batchSize.Add(1)
			return
		}
	}
}

// waitForWrite holds goroutine until batchSize reaches batchMax or until coin flips
func (s *Storage) waitForWrite() error {
	heads := s.coin.Load()
	for {
		if s.coin.Load() != heads {
			err := s.write()
			return err
		}
		if s.batchSize.Load() >= s.batchMax {
			err := s.write()
			return err
		}
	}
}

// write actually writes batch to a file.
func (s *Storage) write() error {
	// ensures only one goroutine will be writing to a file
	for {
		old := s.writeHappens.Load()
		if old {
			return nil
		}
		if s.writeHappens.CompareAndSwap(old, true) {
			break
		}
	}
	// I'm the only one here
	defer s.writeHappens.Store(false)
	defer s.batchSize.Store(0)
	batch := []byte(s.batch.Load())
	batchLen := len(batch)
	n, err := s.writer.Write(batch)
	if err != nil {
		return ErrWalWriteFailed
	}
	// wipe only n written bytes
	if n != batchLen {
		s.batch.Store(string(batch[batchLen+1:]))
	} else {
		s.batch.Store("")
	}
	return writeOk
}

// lockOrWrite stores args to a batch and then either writes them to a file
// or waits until someone else does
func (s *Storage) lockOrWrite(args ...string) error {
	s.addToBuff(args...)
	return s.waitForWrite()
}

// coinFlipper flips a coin every flipTimer interval.
func (s *Storage) coinFlipper(closer chan struct{}) {
	t := time.NewTicker(s.flipTimer)
	defer t.Stop()
	for {
		select {
		case <-closer:
			t.Stop()
			closer <- struct{}{}
			return
		default:
		}

		select {
		case <-t.C:
			for {
				old := s.coin.Load()
				if s.coin.CompareAndSwap(old, !old) {
					break
				}
			}
		case <-closer:
			t.Stop()
			closer <- struct{}{}
			return
		}
	}
}
