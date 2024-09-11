package parser

import (
	"bufio"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"slices"
	"strings"
)

const eol = '\n'
const trim = " \t\n"
const sep = " "

type Parser interface {
	Read() Command
}

type Command struct {
	Command string
	Args    []string
}

type BuffParser struct {
	reader *bufio.Reader
}

// New creates new buffer to read input from
func (bp *BuffParser) New(in io.Reader) {
	bp.reader = bufio.NewReader(in)
}

// Read reads buffer input and tries to compose it into valid Command struct
func (bp *BuffParser) Read() (Command, error) {
	in, err := bp.reader.ReadString(eol)
	if err != nil && err != io.EOF {
		return Command{}, fmt.Errorf("failed to read command: %w", err)
	}

	r, err := composeCommand(strings.Trim(in, trim))
	if err != nil {
		return Command{}, fmt.Errorf("parsing error: %w", err)
	}

	return r, nil
}

// trimArgs composes slice with only args present
func trimArgs(s string) []string {
	// dunno how to parametrize \t
	s = strings.ReplaceAll(s, "\t", sep)
	arr := strings.Split(s, sep)
	arr = slices.DeleteFunc(arr, func(s string) bool {
		return s == ""
	})

	return arr
}

// composeCommand returns valid Command struct
func composeCommand(s string) (Command, error) {
	arr := trimArgs(s)
	err := validateArgs(arr)
	if err != nil {
		return Command{}, fmt.Errorf("argument validation error: %w", err)
	}

	return Command{Command: arr[0], Args: arr[1:]}, nil
}

// validateArgs ensures only correct values are present in the input
func validateArgs(s []string) error {
	validCommands := []string{"GET", "SET", "DEL"}
	ln := len(s)
	val := validator.New(validator.WithRequiredStructEnabled())
	tag := "printascii,containsany=*_/|alphanum|numeric|alpha"

	// min arg length
	if ln < 2 {
		return fmt.Errorf("expected at least 2 arguments, got %d", ln)
	}

	// command is valid
	if !slices.Contains(validCommands, s[0]) {
		return fmt.Errorf("invalid command: %s", s[0])
	}

	// commands have necessary args
	switch s[0] {
	case "GET":
		if ln != 2 {
			return fmt.Errorf("expected 2 arguments, got %d", ln)
		}
	case "DEL":
		if ln != 2 {
			return fmt.Errorf("expected 2 arguments, got %d", ln)
		}
	case "SET":
		if ln != 3 {
			return fmt.Errorf("expected 3 arguments, got %d", ln)
		}
	}

	// validate args
	for i := 1; i < ln; i++ {
		err := val.Var(s[i], tag)
		if err != nil {
			return fmt.Errorf("invalid argument %d: expected %s", i+1, tag)
		}
	}

	return nil
}
