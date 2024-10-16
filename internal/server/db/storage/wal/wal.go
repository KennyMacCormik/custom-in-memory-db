package wal

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/parser"
	"errors"
	"fmt"
	atomicUber "go.uber.org/atomic"
	"log/slog"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ErrWalWriteFailed is the error returned by w.write when batch failed to write to wal.
// (w.write must return EOF itself, not an error wrapping EOF, because callers will test for EOF using ==.)
var ErrWalWriteFailed = errors.New("wal write failed")
var writeOk = errors.New("ok")

// Storage is the same as map.Storage and adds wal implementation.
// Expect files to be written to a WAL_SEG_PATH
type Storage struct {
	initDone bool

	mapMtx sync.Mutex
	m      map[string]string

	batch     atomicUber.String
	batchMax  int32
	batchSize atomic.Int32
	// Timer implementation.
	// We're flipping a coin every flipTimer duration
	flipTimer time.Duration
	coin      atomic.Bool
	closer    chan struct{}
	// writeHappens replaces mutex
	writeHappens atomic.Bool
	writer       writer
}

// Set sets provided value for the provided key.
// Set is thread-safe
func (s *Storage) Set(key, value string) error {
	err := s.lockOrWrite("SET", key, value)
	if err != nil && err != writeOk {
		return err
	}
	return s.set(key, value)
}

// Del removes provided key and it's value.
// Del is thread-safe
func (s *Storage) Del(key string) error {
	err := s.lockOrWrite("DEL", key)
	if err != nil && err != writeOk {
		return err
	}
	return s.del(key)
}

// Get returns a value of the provided key.
// Get is thread-safe
func (s *Storage) Get(key string) (string, error) {
	s.mapMtx.Lock()
	val, ok := s.m[key]
	s.mapMtx.Unlock()
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}

	return val, nil
}

// New used to initialize Storage.
// Any initializations after the first one won't take effect
func (s *Storage) New(conf cmd.Config) error {
	const suf = "WalStorage.New()"
	if !s.initDone {
		s.m = make(map[string]string)
		s.batchMax = int32(conf.Wal.BatchMax)
		s.flipTimer = conf.Wal.BatchTimeout
		s.writeHappens.Store(false)
		s.closer = make(chan struct{})
		go s.coinFlipper(s.closer)
		s.initDone = true
		err := s.writer.New(conf)
		if err != nil {
			return fmt.Errorf("%s failed: %w", suf, err)
		}
	}

	return nil
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

func (s *Storage) Recover(conf cmd.Config, pr parser.Parser, lg *slog.Logger) error {
	return s.writer.Recover(conf, s.set, s.del, pr, lg)
}

// addToBuff loads command to batch and increments batchSize.
func (s *Storage) addToBuff(args ...string) {
	cmnd := strings.Join(args[:], " ")
	cmnd = strings.Join([]string{cmnd, "\n"}, "")

	for {
		old := s.batch.Load()
		new := strings.Join([]string{old, cmnd}, "")
		if !s.writeHappens.Load() && s.batchSize.Load() < s.batchMax && s.batch.CompareAndSwap(old, new) {
			// This doesn't ensure strict batch size. Overflow might happen
			// Is it possible to make two CAS at once?
			s.batchSize.Add(1)
			return
		}
	}
}

// waitForWrite holds goroutine until batchSize reaches batchMax or until coin flips
func (s *Storage) waitForWrite() error {
	heads := s.coin.Load()
	for {
		if s.coin.Load() != heads || s.batchSize.Load() >= s.batchMax {
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
	// TODO handle error
	batch := []byte(s.batch.Load())
	batchLen := len(batch)
	n, err := s.writer.Write(batch)
	// wipe only n written bytes
	if n != batchLen {
		s.batch.Store(string(batch[batchLen+1:]))
	} else {
		s.batch.Store("")
	}
	if err != nil {
		return ErrWalWriteFailed
	}
	return writeOk
}

// lockOrWrite stores args to a batch and then either writes them to a file
// or waits until someone else does
func (s *Storage) lockOrWrite(args ...string) error {
	s.addToBuff(args...)
	return s.waitForWrite()
}

func (s *Storage) set(key, value string) error {
	s.mapMtx.Lock()
	defer s.mapMtx.Unlock()
	s.m[key] = value

	return nil
}

func (s *Storage) del(key string) error {
	s.mapMtx.Lock()
	defer s.mapMtx.Unlock()
	_, ok := s.m[key]
	if !ok {
		return fmt.Errorf("key %s not found", key)
	}
	delete(s.m, key)

	return nil
}

// coinFlipper flips a coin every flipTimer interval.
func (s *Storage) coinFlipper(closer chan struct{}) {
	t := time.NewTicker(s.flipTimer)
	defer t.Stop()
	for {
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
