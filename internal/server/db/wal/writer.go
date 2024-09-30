package wal

import (
	"custom-in-memory-db/internal/server/cmd"
	"errors"
	"fmt"
	"os"
	"strconv"
)

// errMargin represents the amount of bytes we add to our calculations
// to avoid segment files grow more than WAL_SEG_SIZE
const errMargin = 10

// writer actually writes data on disk.
// It is responsible for rotating wal segments
type writer struct {
	// wal file
	walDir     string
	walMaxSize int64
	currSeg    int
	currFile   *os.File
}

// Close gracefully closes currently opened file
func (w *writer) Close() error {
	return w.currFile.Close()
}

// New searches WAL_SEG_PATH for files with integers in their names
// and finds out if we can continue to write to the last segment,
// or it is time to rotate
func (w *writer) New(conf cmd.Config) error {
	w.walDir = conf.Wal.WAL_SEG_PATH
	w.walMaxSize = int64(conf.Wal.WAL_SEG_SIZE)

	if err := w.getCurrSeg(conf); err != nil {
		return fmt.Errorf("getCurrSeg failed: %w", err)
	}

	if !conf.Wal.WAL_SEG_RECOVER {
		if err := w.newSegment(); err != nil {
			return fmt.Errorf("newSegment failed: %w", err)
		}
	} else {
		if err := w.recover(); err != nil {
			return fmt.Errorf("recover failed: %w", err)
		}
	}

	return nil
}

// WriteAndRotate writes wal files and ensures no segment file never grows above WAL_SEG_SIZE
func (w *writer) WriteAndRotate(s []byte) error {
	if !w.isRotate(s) {
		err := w.tryWrite(s)
		if err != nil {
			return fmt.Errorf("tryWrite failed: %w", err)
		}
		return nil
	}

	if !w.isOverflow(s) {
		err := w.tryWrite(s)
		if err != nil {
			return fmt.Errorf("tryWrite failed: %w", err)
		}

		err = w.tryRotate()
		if err != nil {
			return fmt.Errorf("tryRotate failed: %w", err)
		}
		return nil
	}

	index := w.getRotationIndex(s)
	err := w.tryWrite(s[:index])
	if err != nil {
		return fmt.Errorf("tryWrite failed: %w", err)
	}
	err = w.tryRotate()
	if err != nil {
		return fmt.Errorf("tryRotate failed: %w", err)
	}
	err = w.tryWrite(s[index:])
	if err != nil {
		return fmt.Errorf("tryWrite failed: %w", err)
	}

	return nil
}

// getRotationIndex finds out how much we can write to currFile
// without exceeding WAL_SEG_SIZE and all commands intact
func (w *writer) getRotationIndex(s []byte) int {
	st, err := os.Stat(w.walDir + strconv.Itoa(w.currSeg))
	if err != nil {
		return 1
	}

	leftInFile := w.walMaxSize - st.Size()
	for leftInFile > -1 {
		if s[leftInFile] == '\n' {
			break
		}
		leftInFile--
	}

	return int(leftInFile + 1)
}

// isRotate defines if seg file needs rotation after writing s bytes.
// Uses errMargin in its calculations
func (w *writer) isRotate(s []byte) bool {
	st, err := os.Stat(w.walDir + strconv.Itoa(w.currSeg))
	if err != nil {
		return false
	}

	return st.Size()+int64(len(s))+errMargin >= w.walMaxSize
}

// isOverflow defines if seg file will exceed WAL_SEG_SIZE after writing s bytes.
func (w *writer) isOverflow(s []byte) bool {
	st, err := os.Stat(w.walDir + strconv.Itoa(w.currSeg))
	if err != nil {
		return false
	}

	return st.Size()+int64(len(s))-w.walMaxSize > 0
}

// tryWrite writes s bytes to currFile and calls fsync
func (w *writer) tryWrite(s []byte) error {
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

// tryRotate closes tryRotate and creates new seg file
func (w *writer) tryRotate() error {
	_, err := os.Stat(w.walDir + strconv.Itoa(w.currSeg))
	if err != nil {
		return fmt.Errorf("os.Stat failed: %w", err)
	}

	err = w.currFile.Close()
	if err != nil {
		return fmt.Errorf("cannot close wal file: %w", err)
	}

	err = w.newSegment()
	if err != nil {
		return fmt.Errorf("newSegment failed: %w", err)
	}

	return nil
}

// newSegment creates new segment file if no file named currSeg exists.
// Opens existing seg file otherwise
func (w *writer) newSegment() error {
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
	file, err := os.OpenFile(w.walDir+strconv.Itoa(w.currSeg), os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("os.Open failed: %v", err)
	}
	w.currFile = file
	return nil
}

// getCurrSeg finds segment with max number and calculates
// if we will write to it, or rotate to a new one
func (w *writer) getCurrSeg(conf cmd.Config) error {
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
	if st.Size()+errMargin <= w.walMaxSize {
		// means we will write to the current one
		w.currSeg--
		return nil
	}

	return nil
}

// recover reads wal files and runs every command within
func (w *writer) recover() error {
	//TODO implement
	return nil
}
