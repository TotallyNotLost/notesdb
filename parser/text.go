package parser

import (
	"github.com/TotallyNotLost/notesdb/entry"
)

type textParser struct{}

func (p textParser) canParse(source string) bool { return true }
func (p textParser) parse(text string, e *entry.Entry) error {
	e.Type = entry.EntryTypeText
	return nil
}
