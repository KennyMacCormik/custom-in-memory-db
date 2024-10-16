package wal

import (
	"bufio"
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/parser"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"slices"
	"strconv"
)

// errMargin represents the number of bytes we add to our calculations
// to avoid segment files grow more than WAL_SEG_SIZE
const errMargin = 10

// writer implements io.Writer interface.
// it is responsible for rotating wal segments.
// the writer treats any file with natural number as its filename as a wal segment file.
// Wal order assumed to be reflected by the order of the segment files
type writer struct {
	initDone bool

	walDir     string
	walMaxSize int64
	currSeg    int
	currFile   *os.File
}

// Recover loads wal to the running storage.Storage
func (w *writer) Recover(conf cmd.Config, setFunc func(k, v string) error, delFunc func(k string) error, pr parser.Parser, lg *slog.Logger) error {
	const suf = "MapStorage.writer.Recover()"
	files, _, err := w.getFiles(conf, false, true)
	if err != nil {
		return fmt.Errorf("%s.getFiles() failed: %w", suf, err)
	}

	for _, file := range files {
		f, err := os.Open(file.Name())
		if err != nil {
			return fmt.Errorf("%s.os.Open() failed: %w", suf, err)
		}

		w.loadFile(f, setFunc, delFunc, pr, lg)
	}

	return nil
}

// Close gracefully stops writer
func (w *writer) Close() error {
	return w.currFile.Close()
}

// New searches WAL_SEG_PATH for files with natural numbers as their names
// and finds out if we can continue to write to the last one, or it is time to rotate.
// Any initializations after the first one won't take effect
func (w *writer) New(conf cmd.Config) error {
	if !w.initDone {
		w.walDir = conf.Wal.SegPath
		w.walMaxSize = int64(conf.Wal.SegSize)

		if err := w.getCurrSeg(conf); err != nil {
			return fmt.Errorf("getCurrSeg failed: %w", err)
		}

		if err := w.newSegment(); err != nil {
			return fmt.Errorf("newSegment failed: %w", err)
		}

		w.initDone = true
	}

	return nil
}

// WriteAndRotate writes wal files and ensures no segment file ever grows above WAL_SEG_SIZE
func (w *writer) Write(s []byte) (int, error) {
	if !w.isRotate(s) {
		return w.write(s)
	}

	if !w.isOverflow(s) {
		n, err := w.write(s)
		return w.rotate(n, err)
	}

	index := w.getRotationIndex(s)
	n, err := w.write(s[:index])
	n, err = w.rotate(n, err)
	if err != nil {
		return n, err
	}
	return w.write(s)
}

// write decorates tryWrite
func (w *writer) write(s []byte) (int, error) {
	n, err := w.tryWrite(s)
	if err != nil {
		return -1, fmt.Errorf("tryWrite failed: %w", err)
	}
	return n, nil
}

// rotate decorates tryRotate
func (w *writer) rotate(n int, err error) (int, error) {
	if err == nil {
		err = w.tryRotate()
		if err != nil {
			return -1, fmt.Errorf("tryRotate failed: %w", err)
		}
		return n, nil
	}

	return -1, err
}

// loadFile reads commands from a provided file and commits them to the running storage.Storage
func (w *writer) loadFile(f *os.File, set func(k, v string) error, del func(k string) error, pr parser.Parser, lg *slog.Logger) {
	reader := bufio.NewReader(f)
	var err error = nil
	for ; err == nil; _, err = reader.Peek(1) {
		c, e := pr.Read(reader, lg)
		if e != nil {
			// skip incorrect line in a file
			continue
		}
		if c.Command == "SET" {
			_ = set(c.Arg1, c.Arg2)
			continue
		}
		_ = del(c.Arg1)
	}
}

// getRotationIndex finds out how much we can write to currFile
// without exceeding WAL_SEG_SIZE and all commands intact
func (w *writer) getRotationIndex(s []byte) int {
	st, err := os.Stat(path.Join(w.walDir, strconv.Itoa(w.currSeg)))
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

// isRotate defines if a seg file needs rotation after writing s bytes.
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
func (w *writer) tryWrite(s []byte) (int, error) {
	// ensure file exists
	_, err := os.Stat(w.walDir + strconv.Itoa(w.currSeg))
	if err != nil {
		return -1, fmt.Errorf("os.Stat failed: %w", err)
	}
	// write data
	n, err := w.currFile.Write(s)
	if err != nil {
		return -1, fmt.Errorf("write failed: %w", err)
	}
	// fsync
	err = w.currFile.Sync()
	if err != nil {
		return -1, fmt.Errorf("sync failed: %w", err)
	}

	return n, nil
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

// newSegment creates a new segment file if no currSeg exists. Opens an existing seg file otherwise
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
	// open an existing file
	file, err := os.OpenFile(w.walDir+strconv.Itoa(w.currSeg), os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("os.Open failed: %v", err)
	}
	w.currFile = file
	return nil
}

// getCurrSeg finds a segment with max number and calculates
// if we write to it, or rotate to a new one
func (w *writer) getCurrSeg(conf cmd.Config) error {
	files, maxIndex, err := w.getFiles(conf, true, false)
	if err != nil {
		return fmt.Errorf("os.ReadDir failed: %w", err)
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

func (w *writer) getFiles(conf cmd.Config, MaxIndex bool, Sort bool) ([]os.DirEntry, int, error) {
	var maxIndex int
	// list files in folder
	files, err := os.ReadDir(conf.Wal.SegPath)
	if err != nil {
		return nil, 0, fmt.Errorf("os.ReadDir failed: %w", err)
	}
	// filter out non-integers
	for i := 0; i < len(files); i++ {
		_, err := strconv.Atoi(files[i].Name())
		if err != nil {
			files = append(files[:i], files[i+1:]...)
			i--
		}
	}
	// find wal file with max index
	if MaxIndex {
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
	}
	// sorts files in ascending order
	if Sort {
		slices.SortFunc(files, func(a, b os.DirEntry) int {
			n1, _ := strconv.Atoi(a.Name())
			n2, _ := strconv.Atoi(b.Name())
			if n1 < n2 {
				return -1
			}
			if n1 > n2 {
				return 1
			}
			return 0
		})
	}

	return files, maxIndex, nil
}
