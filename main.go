package main

import (
	"flag"
	"fmt"
	"os"
	// "regexp"
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
		fmt.Fprintln(os.Stderr, "Invalid argument: %s", f)
		os.Exit(1)
	}

	fmt.Printf("r = %v\n", *r)
	fmt.Printf("f = %v\n", *f)
	fmt.Printf("fields = %v\n", fields)
}
