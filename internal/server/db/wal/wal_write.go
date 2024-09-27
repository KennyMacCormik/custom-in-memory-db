package wal

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func (w *Wal) write(s []byte) error {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	err := w.tryWrite(s)
	if err != nil {
		return fmt.Errorf("tryWrite failed: %w", err)
	}

	// is it worth to unlock mutex here and lock again to allow more writes to happen?
	err = w.tryRotate()
	if err != nil {
		return fmt.Errorf("tryRotate failed: %w", err)
	}

	return nil
}

func (w *Wal) tryWrite(s []byte) error {
	// ensure file exists
	_, err := os.Stat(w.walDir + strconv.Itoa(w.currSeg))
	if err != nil {
		return fmt.Errorf("os.Stat failed: %w", err)
	}
	// write data
	_, err = w.currFile.Write(s)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	// fsync
	err = w.currFile.Sync()
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	return nil
}

func (w *Wal) tryRotate() error {
	st, err := os.Stat(w.walDir + strconv.Itoa(w.currSeg))
	if err != nil {
		return fmt.Errorf("os.Stat failed: %w", err)
	}

	if st.Size() > w.walMaxSize {
		w.currFile.Close()
		err = w.newSegment()
		if err != nil {
			return fmt.Errorf("newSegment failed: %w", err)
		}
	}

	return nil
}

func (w *Wal) newSegment() error {
	w.currSeg++
	// no next file exists
	if _, err := os.Stat(w.walDir + strconv.Itoa(w.currSeg)); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(w.walDir + strconv.Itoa(w.currSeg))
		if err != nil {
			return fmt.Errorf("os.Create failed: %v", err)
		}
		w.currFile = file
		return nil
	}
	// open existing file
	file, err := os.Open(w.walDir + strconv.Itoa(w.currSeg))
	if err != nil {
		return fmt.Errorf("os.Open failed: %v", err)
	}
	w.currFile = file
	return nil
}
