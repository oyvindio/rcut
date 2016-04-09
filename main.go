package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var startEndRe = regexp.MustCompile("^([0-9]+)-([0-9]+)$")
var startRe = regexp.MustCompile("^([0-9]+)-$")
var endRe = regexp.MustCompile("^-([0-9+])$")

type OutputConfig struct {
	Fields                   map[uint]bool
	HasUnboundedStartingFrom bool
	UnboundedStartingFrom    uint
	OnlyDelimited            bool
	OutputDelimiter          string
	ZeroTerminated           bool
}

func (oc OutputConfig) ShouldOutputField(field uint) bool {
	if oc.HasUnboundedStartingFrom && field >= oc.UnboundedStartingFrom {
		return true
	} else {
		return oc.Fields[field]
	}
}

func (oc OutputConfig) GetLineTerminator() string {
	if oc.ZeroTerminated {
		return "\u0000"
	} else {
		return "\n"
	}
}

func createOutputConfig(f string, onlyDelimited bool, outputDelimiter string,
	zeroTerminated bool) (OutputConfig, error) {
	args := strings.Split(strings.TrimSpace(f), ",")
	var oc = OutputConfig{Fields: make(map[uint]bool), OnlyDelimited: onlyDelimited,
		OutputDelimiter: outputDelimiter, ZeroTerminated: zeroTerminated}

	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		switch {
		case startRe.MatchString(arg):
			matches := startRe.FindStringSubmatch(arg)
			start, err := strconv.ParseUint(matches[1], 10, 32)
			if err != nil {
				return oc, err
			}
			oc.HasUnboundedStartingFrom = true
			oc.UnboundedStartingFrom = uint(start)

		case startEndRe.MatchString(arg):
			matches := startEndRe.FindStringSubmatch(arg)
			start, err := strconv.ParseUint(matches[1], 10, 32)
			if err != nil {
				return oc, err
			}
			end, err := strconv.ParseUint(matches[2], 10, 32)
			if err != nil {
				return oc, err
			}

			for i := uint(start); i <= uint(end); i++ {
				oc.Fields[i] = true
			}

		case endRe.MatchString(arg):
			matches := endRe.FindStringSubmatch(arg)
			end, err := strconv.ParseUint(matches[1], 10, 32)
			if err != nil {
				return oc, err
			}
			for i := uint(0); i <= uint(end); i++ {
				oc.Fields[i] = true
			}

		default:
			field, err := strconv.ParseUint(arg, 10, 32)
			if err != nil {
				return oc, err
			}
			oc.Fields[uint(field)] = true
		}
	}
	return oc, nil
}

func printFieldsFromReader(reader io.Reader, oc OutputConfig, re *regexp.Regexp) {
	lineTerminator := oc.GetLineTerminator()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		tokens := re.Split(scanner.Text(), -1)

		if len(tokens) == 1 && oc.OnlyDelimited {
			continue
		}

		output := make([]string, 0)
		for i, token := range tokens {
			if oc.ShouldOutputField(uint(i + 1)) {
				output = append(output, token)
			}
		}
		fmt.Print(strings.Join(output, oc.OutputDelimiter) + lineTerminator)
	}

	if err := scanner.Err(); err != nil {
		die("Could not read: %v\n", err)
	}
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	os.Exit(1)
}

var regex = flag.String("regex", `\s+`, "Regex to split lines on.")
var fields = flag.String("fields", "1", "Field(s) to output. Should be a comma-separated list of the following:\n"+
	" N    N'th field, counted from 1\n"+
	" N-   from N'th field to end of line\n"+
	" N-M  from N'th to M'th (included) field\n"+
	" -M   from first to M'th (included) field\n")
var onlyDelimited = flag.Bool("only-delimited", false,
	"Do not print lines that do not contain the field separator character.")
var outputDelimiter = flag.String("output-delimiter", " ", "Delimiter to use when outputting fields.")
var zeroTerminated = flag.Bool("zero-terminated", false,
	"Terminate output items with a NUL byte rather than a newline.")
var filenames []string

func init() {
	flag.StringVar(regex, "r", `\s+`, "Regex to split lines on.")
	flag.StringVar(fields, "f", "1", "Field(s) to output.")
	flag.BoolVar(onlyDelimited, "o", false,
		"Do not print lines that do not contain the field separator character.")
	flag.BoolVar(zeroTerminated, "z", false, "Terminate output items with a NUL byte rather than a newline.")
	flag.Parse()
	filenames = flag.Args()
}

func main() {
	oc, err := createOutputConfig(*fields, *onlyDelimited, *outputDelimiter, *zeroTerminated)
	if err != nil {
		die("Invalid field(s): %q\n", *fields)
	}

	re, err := regexp.Compile(*regex)
	if err != nil {
		die("Invalid regex: %q\n", *regex)
	}

	if len(filenames) == 0 || (len(filenames) == 1 && filenames[0] == "-") {
		printFieldsFromReader(os.Stdin, oc, re)
	} else {
		for _, filename := range filenames {
			reader, err := os.Open(filename)
			if err != nil {
				die("Could not open %v\n", filename, err)
			}
			printFieldsFromReader(reader, oc, re)
		}
	}
}
