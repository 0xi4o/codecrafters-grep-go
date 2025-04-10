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
		fmt.Printf("exit code: %d\n", 1)
		os.Exit(1)
	}

	fmt.Printf("exit code: %d\n", 0)
	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	alphabets := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	underscore := "_"

	// alphanumeric
	search := strings.ReplaceAll(pattern, "\\w", fmt.Sprintf("\\d%s%s", alphabets, underscore))
	// numeric
	search = strings.ReplaceAll(search, "\\d", numbers)
	// positive group
	if strings.HasPrefix(search, "[") && strings.HasSuffix(search, "]") {
		search = strings.Replace(search, "[", "", 1)
		search = strings.Replace(search, "]", "", 1)
		// negative group
		if strings.HasPrefix(search, "^") {
			search = strings.Replace(search, "^", "", 1)
			w := fmt.Sprintf("%s%s%s", numbers, alphabets, underscore)
			for _, r := range search {
				w = strings.Replace(w, string(r), "", 1)
			}
			search = w
		}
	}

	return bytes.ContainsAny(line, search), nil
}
