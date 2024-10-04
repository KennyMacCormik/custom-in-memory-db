package wal

import (
	"custom-in-memory-db/internal/server/cmd"
	"fmt"
	"time"
)

// TODO actually implement Err chan

// Input represents data from the request.
// NotifyDone is listened by the requester
// and expects struct{} as a signal of success.
// Data contains string representation of the initial command
type Input struct {
	NotifyDone chan struct{}
	Data       []byte
}

// barrier represents the goroutine that intercepts mutating requests.
// It postpones every request until requests reach maxReqWaiting mark
// or until ticker calls. Whichever comes first, barrier writes
// buffer on disk and then signals to the requesters using Input.NotifyDone chan.
// Err chan used to pass irrecoverable errors such as write failures
type barrier struct {
	Err            chan error
	buffer         []byte
	reqWaiting     []chan struct{}
	in             chan Input
	tickerDuration time.Duration
	ticker         *time.Ticker
	done           chan struct{}
	maxReqWaiting  int

	w writer
}

// New initialize barrier with necessary constants and
// initializes writer object, that handles actual disk writes
func (b *barrier) New(conf cmd.Config, w writer) (chan Input, error) {
	b.buffer = make([]byte, 0, conf.Wal.WAL_SEG_SIZE)
	b.reqWaiting = make([]chan struct{}, 0, conf.Net.NET_MAX_CONN)
	b.in = make(chan Input, conf.Net.NET_MAX_CONN)
	b.maxReqWaiting = conf.Wal.WAL_BATCH_SIZE
	b.ticker = time.NewTicker(conf.Wal.WAL_BATCH_TIMEOUT)
	b.tickerDuration = conf.Wal.WAL_BATCH_TIMEOUT
	b.done = make(chan struct{})
	b.Err = make(chan error, 1)
	b.w = w

	return b.in, nil
}

// Close gracefully shuts down barrier
func (b *barrier) Close() error {
	b.done <- struct{}{}
	b.ticker.Stop()
	close(b.in)
	for val := range b.in {
		b.saveInput(val)
	}
	b.writeWal()
	b.w.Close()
	close(b.done)
	close(b.Err)
	// how to pass two errors?
	return nil
}

// Start runs the actual goroutine
func (b *barrier) Start() {
	const suf = "barrier.start()"
	for {
		select {
		case data := <-b.in:
			if len(b.reqWaiting) < b.maxReqWaiting {
				b.saveInput(data)
			} else {
				// write wal on WAL_BATCH_SIZE
				err := b.writeWal()
				if err != nil {
					b.Err <- fmt.Errorf("%s 'case data' failed: %w", suf, err)
					return
				}
				b.buffer = b.buffer[:0]
				b.ticker.Reset(b.tickerDuration)
				b.saveInput(data)
			}
		case <-b.ticker.C:
			// write wal on WAL_BATCH_TIMEOUT
			err := b.writeWal()
			if err != nil {
				b.Err <- fmt.Errorf("%s 'case ticker' failed: %w", suf, err)
				return
			}
			b.buffer = b.buffer[:0]
		case <-b.done:
			// someone called Close()
			return
		}
	}
}

// writeWal calls w to actually write wal
func (b *barrier) writeWal() error {
	const suf = "barrier.writeWal()"

	if err := b.w.WriteAndRotate(b.buffer); err != nil {
		return fmt.Errorf("%s.writer.WriteAndRotate() failed: %w", suf, err)
	}

	b.notifyDone()

	return nil
}

// notifyDone notifies all postponed connections
func (b *barrier) notifyDone() {
	for _, reqDone := range b.reqWaiting {
		reqDone <- struct{}{}
	}
	b.reqWaiting = b.reqWaiting[:0]
}

// saveInput writes Input object to the barrier
func (b *barrier) saveInput(in Input) {
	b.buffer = append(b.buffer, in.Data...)
	b.reqWaiting = append(b.reqWaiting, in.NotifyDone)
}
