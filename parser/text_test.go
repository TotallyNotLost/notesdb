package parser

import (
	"testing"

	"github.com/TotallyNotLost/notesdb/entry"
)

func TestParseText(t *testing.T) {
	parser := &textParser{}

	t.Run("Empty string sets entry type", func(t *testing.T) {
		var e entry.Entry
		parser.parse("", &e)
		if e.Type != entry.EntryTypeText {
			t.Errorf("want %v got %v", entry.EntryTypeText, e.Type)
		}
	})
}
