package scanner

import (
	"bufio"
	"io"
	"strings"
)

func New(reader io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := strings.Index(string(data), "\n---\n"); i >= 0 {
			return i + len("\n---\n"), data[:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	return scanner
}
