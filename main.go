package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var startEndRe = regexp.MustCompile("^([0-9]+)-([0-9]+)$")
var startRe = regexp.MustCompile("^([0-9]+)-$")

func parseFieldArgs(f string) (map[uint]bool, error) {
	args := strings.Split(strings.TrimSpace(f), ",")
	var ranges = make(map[uint]bool)
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if startRe.MatchString(arg) {
			matches := startRe.FindStringSubmatch(arg)
			start, err := strconv.ParseUint(matches[1], 10, 32)
			if err != nil {
				return ranges, err
			}
			ranges[uint(start)] = true
		} else if startEndRe.MatchString(arg) {
			matches := startEndRe.FindStringSubmatch(arg)
			start, err := strconv.ParseUint(matches[1], 10, 32)
			if err != nil {
				return ranges, err
			}
			end, err := strconv.ParseUint(matches[2], 10, 32)
			if err != nil {
				return ranges, err
			}

			for i := uint(start); i <= uint(end); i++ {
				ranges[i] = true
			}

		} else {
			field, err := strconv.ParseUint(arg, 10, 32)
			if err != nil {
				return ranges, err
			}
			ranges[uint(field)] = true
		}

	}
	return ranges, nil
}

func main() {
	// TODO: long flags http://stackoverflow.com/a/19762274/37208
	r := flag.String("r", "", "Regex to split lines on")
	f := flag.String("f", "1", "Field(s) to output. Default: 1")
	flag.Parse()

	ranges, err := parseFieldArgs(*f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid field(s): %s", f)
		os.Exit(1)
	}

	re, err := regexp.Compile(*r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid regex: %s", r)
		os.Exit(1)
	}

	// TODO: handle file input if stdin is empty and we have a positional arg
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tokens := re.Split(scanner.Text(), -1)

		output := make([]string, 0)
		for i, token := range tokens {
			if ranges[uint(i+1)] {
				output = append(output, token)
			}
		}
		fmt.Println(strings.Join(output, " "))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
