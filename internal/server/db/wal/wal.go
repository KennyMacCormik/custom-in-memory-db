package wal

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/storage"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Wal struct {
	st storage.Storage

	walDir     string
	walMaxSize int64
	currSeg    int

	currFile *os.File

	mtx sync.Mutex
}

func (w *Wal) New(conf cmd.Config, st storage.Storage) error {
	w.st = st

	w.walDir = conf.Wal.WAL_SEG_PATH
	w.walMaxSize = int64(conf.Wal.WAL_SEG_SIZE)
	err := w.getCurrSeg(conf)
	if err != nil {
		return fmt.Errorf("getCurrSeg failed: %w", err)
	}

	if !conf.Wal.WAL_SEG_RECOVER {
		err = w.newSegment()
		if err != nil {
			return fmt.Errorf("newSegment failed: %w", err)
		}
	} else {
		err = w.recover()
		if err != nil {
			return fmt.Errorf("recover failed: %w", err)
		}
	}

	return nil
}

func (w *Wal) Close() {
	w.currFile.Close()
}

func (w *Wal) getCurrSeg(conf cmd.Config) error {
	var maxIndex int
	// list files in folder
	files, err := os.ReadDir(conf.Wal.WAL_SEG_PATH)
	if err != nil {
		return fmt.Errorf("os.ReadDir failed: %w", err)
	}
	// find wal file with max index
	for i, file := range files {
		name, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}
		if !file.IsDir() && w.currSeg < name {
			w.currSeg = name
			maxIndex = i
		}
	}
	// check if it is filled
	st, err := os.Stat(files[maxIndex].Name())
	if err != nil {
		return fmt.Errorf("os.Stat failed: %w", err)
	}
	if st.Size() <= w.walMaxSize {
		w.currSeg--
		return nil
	}

	return nil
}

func (w *Wal) recover() error {

	return nil
}
