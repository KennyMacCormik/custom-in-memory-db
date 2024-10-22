package seg

import (
	"custom-in-memory-db/internal/server/cmd"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strconv"
)

// errMargin represents the number of bytes we add to our calculations
// to avoid segment files grow more than WAL_SEG_SIZE
const errMargin = 10

// Segments implements io.Writer interface.
// It is responsible for rotating wal segments.
// Segments treats any file with natural number as its filename as a wal segment file.
// Wal order assumed to be reflected by the order of the segment files
type Segments struct {
	segPath     string
	segMaxSize  int64
	segFiles    []os.DirEntry
	currSegName int
	currSegFile *os.File
}

func New(conf cmd.Config) (*Segments, error) {
	var err error
	seg := Segments{}
	seg.segPath = conf.Wal.SegPath
	seg.segMaxSize = int64(conf.Wal.SegSize)
	seg.segFiles, err = seg.getFiles()
	if err != nil {
		return nil, fmt.Errorf("get segment files failed: %w", err)
	}
	if err := seg.newSegment(); err != nil {
		return nil, fmt.Errorf("newSegment failed: %w", err)
	}
	return &seg, nil
}

func (s *Segments) ExportSegNames() []string {
	var result = make([]string, 0, len(s.segFiles))
	for _, file := range s.segFiles {
		result = append(result, path.Join(s.segPath, file.Name()))
	}
	return result
}

func (s *Segments) Write(n []byte) (int, error) {
	batchLen := int64(len(n))
	if !s.isRotate(batchLen) {
		return s.write(n)
	}

	if !s.isOverflow(batchLen) {
		nn, err := s.write(n)
		return s.rotate(nn, err)
	}

	index := s.getRotationIndex(n)
	nn, err := s.write(n[:index])
	nn, err = s.rotate(nn, err)
	if err != nil {
		return -1, err
	}
	return s.write(n[index:])
}

func (s *Segments) Close() error {
	return s.currSegFile.Close()
}

// getRotationIndex finds out how much we can write to currFile
// without exceeding WAL_SEG_SIZE and all commands intact
func (s *Segments) getRotationIndex(n []byte) int {
	st, err := os.Stat(path.Join(s.segPath, strconv.Itoa(s.currSegName)))
	if err != nil {
		// we will fail later
		return 1
	}
	// move backwards to rotate without breaking commands
	leftInFile := s.segMaxSize - st.Size()
	for leftInFile > -1 {
		if n[leftInFile] == '\n' {
			break
		}
		leftInFile--
	}

	return int(leftInFile + 1)
}

// rotate decorates tryRotate
func (s *Segments) rotate(n int, err error) (int, error) {
	if err == nil {
		err = s.tryRotate()
		if err != nil {
			return -1, err
		}
		return n, nil
	}

	return -1, err
}

// tryRotate closes tryRotate and creates new seg file
func (s *Segments) tryRotate() error {
	// ensure file exists
	pth := path.Join(s.segPath, strconv.Itoa(s.currSegName))
	_, err := os.Stat(pth)
	if err != nil {
		return fmt.Errorf("file %q check failed: %w", pth, err)
	}
	// close current seg file
	err = s.currSegFile.Close()
	if err != nil {
		return fmt.Errorf("cannot close %q file: %w", pth, err)
	}
	// create new seg file
	err = s.newSegment()
	if err != nil {
		return fmt.Errorf("newSegment failed: %w", err)
	}

	return nil
}

// newSegment creates a new segment file if no currSeg exists. Opens an existing seg file otherwise
func (s *Segments) newSegment() error {
	s.currSegName++
	pth := path.Join(s.segPath, strconv.Itoa(s.currSegName))
	// no next file exists
	if _, err := os.Stat(pth); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(pth)
		if err != nil {
			return fmt.Errorf("os.Create %q failed: %v", pth, err)
		}
		s.currSegFile = file
		return nil
	}
	// open an existing file
	file, err := os.OpenFile(pth, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("os.Open failed: %v", err)
	}
	s.currSegFile = file
	return nil
}

// write writes n bytes to currFile and calls fsync
func (s *Segments) write(n []byte) (int, error) {
	// ensure file exists
	pth := path.Join(s.segPath, strconv.Itoa(s.currSegName))
	_, err := os.Stat(pth)
	if err != nil {
		return -1, fmt.Errorf("file %q check failed: %w", pth, err)
	}
	// write data
	nn, err := s.currSegFile.Write(n)
	if err != nil {
		return -1, fmt.Errorf("file %q write failed: %w", pth, err)
	}
	// fsync
	err = s.currSegFile.Sync()
	if err != nil {
		return -1, fmt.Errorf("file %q fsync failed: %w", pth, err)
	}

	return nn, nil
}

// isOverflow defines if seg file will exceed WAL_SEG_SIZE after writing s bytes.
func (s *Segments) isOverflow(num int64) bool {
	st, err := os.Stat(path.Join(s.segPath, strconv.Itoa(s.currSegName)))
	if err != nil {
		return false
	}

	return st.Size()+num-s.segMaxSize > 0
}

// isRotate defines if a seg file needs rotation after writing s bytes.
// Uses errMargin in its calculations
func (s *Segments) isRotate(num int64) bool {
	st, err := os.Stat(path.Join(s.segPath, strconv.Itoa(s.currSegName)))
	if err != nil {
		return false
	}

	return st.Size()+num+errMargin >= s.segMaxSize
}

func (s *Segments) getFiles() ([]os.DirEntry, error) {
	// list files in folder
	tmpFiles, err := os.ReadDir(s.segPath)
	if err != nil {
		return nil, err
	}
	files := make([]os.DirEntry, 0, len(tmpFiles))
	// filter out non-integers
	for _, file := range tmpFiles {
		_, err = strconv.Atoi(file.Name())
		if err == nil {
			files = append(files, file)
		}
	}
	// sort
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
	return files, nil
}

func (s *Segments) getCurrSeg() error {
	files, err := s.getFiles()
	if err != nil {
		return fmt.Errorf("failed reading %q: %w", s.segPath, err)
	}
	if len(files) == 0 {
		return nil
	}
	// check if it is filled
	st, err := os.Stat(files[len(files)-1].Name())
	if err != nil {
		return fmt.Errorf("os.Stat failed: %w", err)
	}
	if st.Size()+errMargin <= s.segMaxSize {
		// means we will write to the current one
		s.currSegName--
		return nil
	}

	return nil
}
