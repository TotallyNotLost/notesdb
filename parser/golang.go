package parser

import (
	"bufio"
	"io"
	"strings"

	"github.com/TotallyNotLost/notesdb/storage"
)

type goParser struct {
	source     string
	scanner    *bufio.Scanner
	textParser *textParser
}

func (p *goParser) SetSource(source string)     { p.textParser.SetSource(source) }
func (p *goParser) SetReader(reader io.Reader)  { p.textParser.SetReader(reader) }
func (p *goParser) canParse(source string) bool { return strings.HasSuffix(string(source), ".go") }
func (p *goParser) Next() (entry storage.Entry, err error) {
	entry, err = p.textParser.Next()
	if err != nil {
		return entry, err
	}
	entry.Type = storage.EntryTypeCode
	return entry, nil
}
