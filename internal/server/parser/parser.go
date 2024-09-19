package parser

import (
	"bufio"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"slices"
	"strings"
)

const Eol = '\n'
const Trim = " \t\n"
const Sep = " "

type Command struct {
	Command string
	Args    []string
}

type Parser interface {
	Read(validCommands []string, lg *slog.Logger) (Command, error)
}

// BufferRead implements reading from provided bufio.Reader.
// Expects bufio.Reader to be initialized
func BufferRead(reader *bufio.Reader, vc []string, lg *slog.Logger) (Command, error) {
	in, err := reader.ReadString(Eol)
	if err != nil && err != io.EOF {
		lg.Error("failed to read command", "error", err.Error())
		return Command{}, fmt.Errorf("failed to read from connection: %w", err)
	}

	lg.Debug("logging cmd", "cmd", in)

	r, err := composeCommand(strings.Trim(in, Trim), vc)
	if err != nil {
		lg.Error("parsing error", "error", err.Error())
		return Command{}, fmt.Errorf("parsing error: %w", err)
	}

	return r, nil
}

// trimArgs composes slice with only args present
func trimArgs(s string) Command {
	// dunno how to parametrize \t
	s = strings.ReplaceAll(s, "\t", Sep)
	arr := strings.Split(s, Sep)
	arr = slices.DeleteFunc(arr, func(s string) bool {
		return s == ""
	})

	return Command{Command: arr[0], Args: arr[1:]}
}

// composeCommand returns valid Command struct
func composeCommand(s string, vc []string) (Command, error) {
	result := trimArgs(s)
	err := validateArgs(result, vc)
	if err != nil {
		return Command{}, fmt.Errorf("argument validation error: %w", err)
	}

	return result, nil
}

// validateArgs ensures only correct values are present in the input
func validateArgs(c Command, vc []string) error {
	ln := len(c.Args)
	val := validator.New(validator.WithRequiredStructEnabled())
	tag := "printascii,containsany=*_/|alphanum|numeric|alpha"

	// command is valid
	if !slices.Contains(vc, c.Command) {
		return fmt.Errorf("invalid command: %s", c.Command)
	}

	// commands have necessary args
	switch c.Command {
	case "GET":
		if ln != 1 {
			return fmt.Errorf("expected 1 argument, got %d", ln)
		}
	case "DEL":
		if ln != 1 {
			return fmt.Errorf("expected 1 argument, got %d", ln)
		}
	case "SET":
		if ln != 2 {
			return fmt.Errorf("expected 2 arguments, got %d", ln)
		}
	case "QUIT", "EXIT":
		if ln != 0 {
			return fmt.Errorf("expected 0 arguments, got %d", ln)
		}
	}

	// validate args
	for i := 0; i < ln; i++ {
		err := val.Var(c.Args[i], tag)
		if err != nil {
			return fmt.Errorf("invalid argument %d: expected %s", i+1, tag)
		}
	}

	return nil
}
