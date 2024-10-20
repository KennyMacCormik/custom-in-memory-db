package parser

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"math/rand"
	"strings"
	"testing"
)

const validSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890*_/"
const num = 100_000
const testArgLen = 20
const unacceptableChars = "!\"#$%&'()+|-.:;<=>?@[]^`{},~"

var nilLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

// get

func TestRead_Get_Positive(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "GET 1\n",
		expected: Command{
			Command: "GET",
			Arg1:    "1",
		},
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Get_Negative_ZeroArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "GET\n",
		expected: Command{},
		err:      "parser.Read().composeCommand().validateArgs() failed: \"GET\" expects exactly 1 arg",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Get_Negative_ExcessiveArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "GET 1 2\n",
		expected: Command{},
		err:      "parser.Read().composeCommand().validateArgs() failed: \"GET\" expects exactly 1 arg",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Get_Negative_NoEndline(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "GET 1",
		expected: Command{},
		err:      "parser.Read() failed: expect '\\n' as EOL, got none",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

// SET

func TestRead_Set_Positive(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "SET 1 2\n",
		expected: Command{
			Command: "SET",
			Arg1:    "1",
			Arg2:    "2",
		},
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Set_Negative_NoArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "SET\n",
		expected: Command{},
		err:      "parser.Read().composeCommand().validateArgs() failed: \"SET\" expects exactly 2 args",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Set_Negative_InsufficientArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "SET 1\n",
		expected: Command{},
		err:      "parser.Read().composeCommand().validateArgs() failed: \"SET\" expects exactly 2 args",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Set_Negative_ExcessiveArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "SET 1 2 3\n",
		expected: Command{},
		err:      "parser.Read().composeCommand().validateArgs() failed: \"SET\" expects exactly 2 args",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Set_Negative_NoEndline(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "SET 1 2 3",
		expected: Command{},
		err:      "parser.Read() failed: expect '\\n' as EOL, got none",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

// Del

func TestRead_Del_Positive(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "DEL 1\n",
		expected: Command{
			Command: "DEL",
			Arg1:    "1",
		},
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Del_Negative_ZeroArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "DEL\n",
		expected: Command{},
		err:      "parser.Read().composeCommand().validateArgs() failed: \"DEL\" expects exactly 1 arg",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Del_Negative_ExcessiveArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "DEL 1 2\n",
		expected: Command{},
		err:      "parser.Read().composeCommand().validateArgs() failed: \"DEL\" expects exactly 1 arg",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Del_Negative_NoEndline(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "DEL 1",
		expected: Command{},
		err:      "parser.Read() failed: expect '\\n' as EOL, got none",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

// Misc

func TestRead_BogusCommand_WithoutArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "QWERTY\n",
		expected: Command{},
		err:      "parser.Read().composeCommand().validateArgs() failed: got empty or unexpected command \"QWERTY\"",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_BogusCommand_WithArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "HELLO 1 2\n",
		expected: Command{},
		err:      "parser.Read().composeCommand().validateArgs() failed: got empty or unexpected command \"HELLO\"",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_OneSymbolNewLine_String(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "\n",
		expected: Command{},
		err:      "parser.Read().composeCommand() failed: empty command",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_ZeroString(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "\n",
		expected: Command{},
		err:      "parser.Read().composeCommand() failed: empty command",
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

// arg validation positive

func TestRead_Positive_alphanum(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "SET q1w3 qwe87sdgbi823948nadf09324h\n",
		expected: Command{
			Command: "SET",
			Arg1:    "q1w3",
			Arg2:    "qwe87sdgbi823948nadf09324h",
		},
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Positive_numeric(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "SET 123830 45612841298\n",
		expected: Command{
			Command: "SET",
			Arg1:    "123830",
			Arg2:    "45612841298",
		},
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Positive_alpha(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "SET qwewfgdfsetjgjol yergEPSBFhgbslkaIBSVDFskdgbjd\n",
		expected: Command{
			Command: "SET",
			Arg1:    "qwewfgdfsetjgjol",
			Arg2:    "yergEPSBFhgbslkaIBSVDFskdgbjd",
		},
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Positive_underscore(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "SET qwewf123gd_fsetj_gj43ol yergEP_SBFhgb46sl_1kaIBSVDF_skdgbjd\n",
		expected: Command{
			Command: "SET",
			Arg1:    "qwewf123gd_fsetj_gj43ol",
			Arg2:    "yergEP_SBFhgb46sl_1kaIBSVDF_skdgbjd",
		},
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Positive_slash(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "SET qwewf123gd/fsetj/gj43ol yergEP/SBFhgb46sl/1kaIBSVDF/skdgbjd\n",
		expected: Command{
			Command: "SET",
			Arg1:    "qwewf123gd/fsetj/gj43ol",
			Arg2:    "yergEP/SBFhgb46sl/1kaIBSVDF/skdgbjd",
		},
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Positive_star(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "SET weEr1yJAR92f3*fs21eDFtj*gj4F3ol 3yer1gE2P*SwB334Fhg5b46sl*1kaIBfS1s23aVDF*sk397dSFgSj2d\n",
		expected: Command{
			Command: "SET",
			Arg1:    "weEr1yJAR92f3*fs21eDFtj*gj4F3ol",
			Arg2:    "3yer1gE2P*SwB334Fhg5b46sl*1kaIBfS1s23aVDF*sk397dSFgSj2d",
		},
	}

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Get_Positive_RandArgs(t *testing.T) {
	var letters = []rune(validSymbols)
	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	pr := New()

	for i := 0; i < num; i++ {
		randStr := randSeq(testArgLen)
		randStr = "GET " + randStr + string(eol)
		val, err := pr.Read(bytes.NewReader([]byte(randStr)), nilLogger)

		assert.NoError(t, err)
		res := strings.Join([]string{val.Command, val.Arg1}, sep)
		res += string(eol)
		assert.Equal(t, res, randStr)
		if err != nil {
			break
		}
	}
}

func TestRead_Del_Positive_RandArgs(t *testing.T) {
	var letters = []rune(validSymbols)
	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	pr := New()

	for i := 0; i < num; i++ {
		randStr := randSeq(testArgLen)
		randStr = "DEL " + randStr + string(eol)
		val, err := pr.Read(bytes.NewReader([]byte(randStr)), nilLogger)

		assert.NoError(t, err)
		res := strings.Join([]string{val.Command, val.Arg1}, sep)
		res += string(eol)
		assert.Equal(t, res, randStr)
		if err != nil {
			break
		}
	}
}

func TestRead_Set_Positive_RandArgs(t *testing.T) {
	var letters = []rune(validSymbols)
	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	pr := New()

	for i := 0; i < num; i++ {
		randStr1 := randSeq(testArgLen)
		randStr2 := randSeq(testArgLen)
		randStr := strings.Join([]string{"SET", randStr1, randStr2}, sep)
		randStr += string(eol)
		val, err := pr.Read(bytes.NewReader([]byte(randStr)), nilLogger)

		assert.NoError(t, err)
		res := strings.Join([]string{val.Command, val.Arg1, val.Arg2}, sep)
		res += string(eol)
		assert.Equal(t, res, randStr)
		if err != nil {
			break
		}
	}
}

// arg validation negative

func TestRead_Negative_Unicode(t *testing.T) {
	arg := "val_⌘"
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		expected: Command{},
	}
	testCase.ioInput = fmt.Sprintf("GET %s\n", arg)
	testCase.err = fmt.Sprintf("parser.Read().composeCommand().validateArgs() failed: got %q, expected %q", arg, tag)

	pr := New()
	r := bytes.NewReader([]byte(testCase.ioInput))

	val, err := pr.Read(r, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_Negative_UnacceptableChars(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		expected: Command{},
	}

	pr := New()

	for _, char := range unacceptableChars {
		arg := fmt.Sprintf("k3y/%s_value*", string(char))
		testCase.ioInput = fmt.Sprintf("SET %s 1\n", arg)
		testCase.err = fmt.Sprintf("parser.Read().composeCommand().validateArgs() failed: got %q, expected %q", arg, tag)

		val, err := pr.Read(bytes.NewReader([]byte(testCase.ioInput)), nilLogger)
		assert.EqualError(t, err, testCase.err)
		assert.Equal(t, testCase.expected, val)
		if err != nil {
			break
		}
	}

	for _, char := range unacceptableChars {
		arg := fmt.Sprintf("k3y/%s_value*", string(char))
		testCase.ioInput = fmt.Sprintf("SET 1 %s\n", arg)
		testCase.err = fmt.Sprintf("parser.Read().composeCommand().validateArgs() failed: got %q, expected %q", arg, tag)

		val, err := pr.Read(bytes.NewReader([]byte(testCase.ioInput)), nilLogger)
		assert.EqualError(t, err, testCase.err)
		assert.Equal(t, testCase.expected, val)
		if err != nil {
			break
		}
	}
}
