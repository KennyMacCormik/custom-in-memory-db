package wal

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/storage"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

// Wal shadows any implementation of storage.Storage.
// It provides an ability to dump any mutating
// commands on disk before commiting them
type Wal struct {
	st storage.Storage

	bar           barrier
	sendToBarrier chan Input

	w WriterInterface
}

func (w *Wal) Recover(conf cmd.Config, lg *slog.Logger) error {
	const suf = "wal.New().barrier.New() failed:"
	// how to unit-test this?
	setter := func(k, v string) error {
		return w.st.Set(k, v)
	}
	deller := func(k string) error {
		return w.st.Del(k)
	}

	err := w.w.Recover(conf, setter, deller, lg)
	if err != nil {
		lg.Error(suf, "error", err.Error())
		return fmt.Errorf("%s: %w", suf, err)
	}

	return nil
}

// New expects initialized storage.Storage object.
// New starts barrier in the separate goroutine.
func (w *Wal) New(conf cmd.Config, st storage.Storage, writer WriterInterface) {
	w.w = writer

	w.st = st
	// New() returns channel where all commands will be sent to
	w.sendToBarrier = w.bar.New(conf, w.w)
}

func (w *Wal) Start() {
	go w.bar.Start()
}

// Close stops barrier and storage.Storage gracefully
func (w *Wal) Close() error {
	// how to pass two errors?
	err1 := w.st.Close()
	err2 := w.bar.Close()
	var errStr string
	if err1 != nil {
		errStr = err1.Error()
	}
	if err2 != nil {
		errStr = strings.Join([]string{err2.Error(), errStr}, "; ")
	}
	if errStr != "" {
		errStr = strings.Join([]string{"wal.Close().storage.Close() and wal.Close().barrier.Close() joined error", errStr}, ": ")
		return errors.New(errStr)
	}
	return nil
}

// Get has no difference with storage.Storage implementation
func (w *Wal) Get(key string) (string, error) {
	const suf = "wal.Get()"
	result, err := w.st.Get(key)
	if err != nil {
		return "", fmt.Errorf("%s : %w", suf, err)
	}

	return result, nil
}

// Set runs waitForWal() func.
// It is the only difference with storage.Storage implementation
func (w *Wal) Set(key, value string) error {
	const suf = "wal.Set()"

	w.waitForWal([]byte("SET " + key + " " + value + "\n"))

	err := w.st.Set(key, value)
	if err != nil {
		return fmt.Errorf("%s : %w", suf, err)
	}

	return nil
}

// Del runs waitForWal() func.
// It is the only difference with storage.Storage implementation
func (w *Wal) Del(key string) error {
	const suf = "wal.Del()"

	w.waitForWal([]byte("DEL " + key + "\n"))

	err := w.st.Del(key)
	if err != nil {
		return fmt.Errorf("%s: %w", suf, err)
	}

	return nil
}

// waitForWal sends command to the barrier and waits until it is commited to wal
func (w *Wal) waitForWal(inData []byte) {
	msg := Input{
		NotifyDone: make(chan struct{}),
		Data:       inData,
	}
	w.sendToBarrier <- msg
	<-msg.NotifyDone
	close(msg.NotifyDone)
}
