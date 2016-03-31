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

func parseFieldArgs(f string) ([]uint64, error) {
	args := strings.Split(strings.TrimSpace(f), ",")
	var fields = make([]uint64, len(args))
	for _, arg := range args {
		field, err := strconv.ParseUint(strings.TrimSpace(arg), 10, 64)
		if err != nil {
			return fields, err
		} else {
			fields = append(fields, field)
		}
	}
	return fields, nil
}

func main() {
	r := flag.String("r", "", "Regex to split lines on")
	f := flag.String("f", "1", "Field(s) to output. Default: 1")
	flag.Parse()

	fields, err := parseFieldArgs(*f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid field(s): %s", f)
		os.Exit(1)
	}

	re, err := regexp.Compile(*r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid regex: %s", r)
		os.Exit(1)
	}

	fmt.Printf("r = %v\n", *r)
	fmt.Printf("f = %v\n", *f)
	fmt.Printf("fields = %v\n", fields)

	// TODO: handle file input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		tokens := re.Split(scanner.Text(), -1)

		output := make([]string, len(tokens))
		for i, token := range tokens {
			if i == 1 {
				output = append(output, token)
			}
		}
		fmt.Println(strings.Join(output, ""))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
