package parser

import (
	"strings"

	"github.com/TotallyNotLost/notesdb/entry"
)

type goParser struct{}

func (p goParser) canParse(source string) bool { return strings.HasSuffix(string(source), ".go") }
func (p goParser) parse(text string, e *entry.Entry) error {
	e.Type = entry.EntryTypeCode
	return nil
}
