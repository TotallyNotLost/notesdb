package parser

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/TotallyNotLost/notesdb/entry"
)

type textParser struct{}

func (p textParser) canParse(source string) bool { return true }
func (p textParser) parse(source string, text string) (e entry.Entry, err error) {
	h := sha1.New()
	h.Write([]byte(text))
	// TODO: Why doesn't this match sha1sum?
	id := hex.EncodeToString(h.Sum(nil))
	e = entry.NewEntry(id, source, entry.EntryTypeText, source, text, []string{})
	return e, nil
}
