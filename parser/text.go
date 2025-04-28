package parser

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"io"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/TotallyNotLost/notesdb/scanner"
)

type textParser struct {
	source  string
	scanner *bufio.Scanner
}

func (p *textParser) SetSource(source string)     { p.source = source }
func (p *textParser) SetReader(reader io.Reader)  { p.scanner = scanner.New(reader) }
func (p *textParser) canParse(source string) bool { return true }
func (p *textParser) Next() (e entry.Entry, err error) {
	more := p.scanner.Scan()
	if !more {
		return e, io.EOF
	}
	body := p.scanner.Text()
	h := sha1.New()
	h.Write([]byte(body))
	// TODO: Why doesn't this match sha1sum?
	id := hex.EncodeToString(h.Sum(nil))
	e = entry.NewEntry(id, p.source, entry.EntryTypeText, p.source, body, []string{})
	return e, nil
}
