package wal

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/storage"
	"fmt"
	"log/slog"
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
func (w *Wal) New(conf cmd.Config, st storage.Storage, wrtr WriterInterface) error {
	const suf = "wal.New().barrier.New()"
	var err error

	w.w = wrtr

	w.st = st
	// New() returns channel where all commands will be sent to
	w.sendToBarrier, err = w.bar.New(conf, w.w)
	if err != nil {
		return fmt.Errorf("%s failed: %w", suf, err)
	}

	return nil
}

func (w *Wal) Start() {
	go w.bar.Start()
}

// Close stops barrier and storage.Storage gracefully
func (w *Wal) Close() error {
	w.st.Close()
	w.bar.Close()
	// how to pass two errors?
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

// waitForWal sends command to barrier and waits until it is commited to wal
func (w *Wal) waitForWal(inData []byte) {
	msg := Input{
		NotifyDone: make(chan struct{}),
		Data:       inData,
	}
	w.sendToBarrier <- msg
	<-msg.NotifyDone
	close(msg.NotifyDone)
}
