package parser

import (
	"bufio"
	"custom-in-memory-db/internal/server/cmd"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"slices"
	"strings"
)

//const Eol = '\n'
//const Trim = " \t\n"
//const Sep = " "
//const ToReplace = "\t"

type Command struct {
	Command string
	Arg1    string
	Arg2    string
}

type Parser interface {
	Read(io.Reader, *slog.Logger) (Command, error)
}

type Parse struct {
	initDone bool

	eol            byte
	trim           string
	sep            string
	toReplaceBySep string
	tag            string
}

// New used to initialize Storage.
// Any initializations after the first one won't take effect
func (p *Parse) New(conf cmd.Config) {
	if !p.initDone {
		p.initDone = true
		p.eol = conf.Parser.Eol
		p.trim = conf.Parser.Trim
		p.sep = conf.Parser.Sep
		p.toReplaceBySep = conf.Parser.ToReplaceBySep
		p.tag = conf.Parser.Tag
	}
}

// Read reads r until Parse.eol and converts it to Command
func (p *Parse) Read(r io.Reader, lg *slog.Logger) (Command, error) {
	const suf = "parser.Read()"

	bufR := bufio.NewReader(r)
	str, err := bufR.ReadString(p.eol)
	if err != nil {
		if err == io.EOF {
			lg.Error(fmt.Sprintf("%s failed: expect %q as EOL, got none", suf, p.eol), "error", err.Error())
			return Command{}, fmt.Errorf("%s failed: expect %q as EOL, got none", suf, p.eol)
		}
		lg.Error(fmt.Sprintf("%s failed", suf), "error", err.Error())
		return Command{}, fmt.Errorf("%s failed: %w", suf, err)
	}
	lg.Debug(suf, "input", str)

	cmnd, err := p.composeCommand(strings.Trim(str, p.trim))
	if err != nil {
		lg.Error(suf, "error", err.Error())
		return Command{}, err
	}

	return cmnd, nil
}

// composeCommand returns valid Command struct
func (p *Parse) composeCommand(s string) (Command, error) {
	if s == "" {
		return Command{}, errors.New("parser.Read().composeCommand() failed: empty command")
	}

	result := p.trimArgs(s)
	err := p.validateArgs(result)
	if err != nil {
		return Command{}, err
	}

	return result, nil
}

// trimArgs composes slice with only args present
func (p *Parse) trimArgs(s string) Command {
	s = strings.ReplaceAll(s, p.toReplaceBySep, p.sep)
	arr := strings.Split(s, p.sep)
	arr = slices.DeleteFunc(arr, func(s string) bool {
		return s == ""
	})

	if len(arr) == 2 {
		return Command{Command: arr[0], Arg1: arr[1]}
	}

	if len(arr) == 3 {
		return Command{Command: arr[0], Arg1: arr[1], Arg2: arr[2]}
	}

	return Command{Command: arr[0], Arg1: "", Arg2: ""}
}

// validateArgs ensures only correct values are present in the input
func (p *Parse) validateArgs(c Command) error {
	const suf = "parser.Read().composeCommand().validateArgs()"
	val := validator.New(validator.WithRequiredStructEnabled())
	validate := func(arg string) error {
		err := val.Var(arg, p.tag)
		if err != nil {
			return fmt.Errorf("%s failed: got %q, expected %q", suf, arg, p.tag)
		}
		return nil
	}

	switch c.Command {
	case "GET":
		if c.Arg2 != "" || c.Arg1 == "" && c.Arg2 == "" {
			return fmt.Errorf("%s failed: %q expects exactly 1 arg", suf, c.Command)
		}
		return validate(c.Arg1)
	case "DEL":
		if c.Arg2 != "" || c.Arg1 == "" && c.Arg2 == "" {
			return fmt.Errorf("%s failed: %q expects exactly 1 arg", suf, c.Command)
		}
		return validate(c.Arg1)
	case "SET":
		if c.Arg2 == "" && c.Arg1 == "" || c.Arg1 == "" || c.Arg2 == "" {
			return fmt.Errorf("%s failed: %q expects exactly 2 args", suf, c.Command)
		}
		err := validate(c.Arg1)
		if err != nil {
			return err
		}
		return validate(c.Arg2)
	default:
		return fmt.Errorf("%s failed: got empty or unexpected command %q", suf, c.Command)
	}
}
