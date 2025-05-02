package parser

import (
	"fmt"
	"strings"

	"github.com/TotallyNotLost/notesdb/entry"
)

type goParser struct{}

func (p goParser) Parse(source string, text string) (e entry.Entry, err error) {
	if !strings.HasSuffix(string(source), ".go") {
		return e, fmt.Errorf("Can't parse")
	}
	e.Type = entry.EntryTypeCode
	return e, nil
}
