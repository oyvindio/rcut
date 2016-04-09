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

var r, f string
var filenames []string
var startEndRe = regexp.MustCompile("^([0-9]+)-([0-9]+)$")
var startRe = regexp.MustCompile("^([0-9]+)-$")
var endRe = regexp.MustCompile("^-([0-9+])$")

type OutputConfig struct {
	Fields                   map[uint]bool
	HasUnboundedStartingFrom bool
	UnboundedStartingFrom    uint
}

func (oc OutputConfig) ShouldOutputField(field uint) bool {
	if oc.HasUnboundedStartingFrom && field >= oc.UnboundedStartingFrom {
		return true
	} else {
		return oc.Fields[field]
	}
}

func createOutputConfig(f string) (OutputConfig, error) {
	args := strings.Split(strings.TrimSpace(f), ",")
	var oc = OutputConfig{Fields: make(map[uint]bool)}
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
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		tokens := re.Split(scanner.Text(), -1)

		output := make([]string, 0)
		for i, token := range tokens {
			if oc.ShouldOutputField(uint(i + 1)) {
				output = append(output, token)
			}
		}
		fmt.Println(strings.Join(output, " "))
	}

	if err := scanner.Err(); err != nil {
		die("Could not read: %v\n", err)
	}
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	os.Exit(1)
}

func init() {
	// TODO: long flags http://stackoverflow.com/a/19762274/37208
	flag.StringVar(&r, "r", `\s+`, "Regex to split lines on.")
	flag.StringVar(&f, "f", "1", "Field(s) to output.")
	flag.Parse()
	filenames = flag.Args()
}

func main() {
	oc, err := createOutputConfig(f)
	if err != nil {
		die("Invalid field(s): %q\n", f)
	}

	re, err := regexp.Compile(r)
	if err != nil {
		die("Invalid regex: %q\n", r)
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
