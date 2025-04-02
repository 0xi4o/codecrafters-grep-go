package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	// alphanumeric
	search := strings.ReplaceAll(pattern, "\\w", "\\dabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")
	// numeric
	search = strings.ReplaceAll(search, "\\d", "0123456789")
	// positive group
	if strings.HasPrefix(search, "[") && strings.HasSuffix(search, "]") {
		search = strings.Replace(search, "[", "", 1)
		search = strings.Replace(search, "]", "", 1)
		// negative group
		if strings.HasPrefix(search, "^") {
			search = strings.Replace(search, "^", "", 1)
			fmt.Printf("negative character group search: %s\n", search)
			return !bytes.ContainsAny(line, search), nil
		}
	}

	return bytes.ContainsAny(line, search), nil
}
