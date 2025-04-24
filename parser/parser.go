package parser

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/TotallyNotLost/notesdb/scanner"
)

type Parser interface {
	SetSource(source string)
	SetReader(io.Reader)
	canParse(source string) bool
	Next() (entry.Entry, error)
}

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

func getMetadata(text string) map[string][]string {
	lines := strings.Split(text, "\n")

	o := make(map[string][]string)

	for _, l := range lines {
		r, _ := regexp.Compile("^\\[_metadata_:*(\\w+)\\]:# \"(.*)\"$")

		if !r.MatchString(l) {
			continue
		}

		key := r.FindStringSubmatch(l)[1]
		if _, ok := o[key]; !ok {
			o[key] = []string{}
		}

		value := r.FindStringSubmatch(l)[2]
		o[key] = append(o[key], value)
	}

	return o
}

func NewRelative(relationship string) (*entry.Relative, error) {
	r := new(entry.Relative)

	kvPairs := strings.SplitSeq(relationship, ",")

	for kvPair := range kvPairs {
		k := strings.Split(kvPair, "=")[0]
		v := strings.Split(kvPair, "=")[1]

		var err error

		switch k {
		case "id":
			r.Id = v
		case "start":
			r.Start, err = strconv.Atoi(v)
			if err != nil {
				return new(entry.Relative), fmt.Errorf("Error parsing relation start: %w", err)
			}
		case "end":
			r.End, err = strconv.Atoi(v)
			if err != nil {
				return new(entry.Relative), fmt.Errorf("Error parsing relation end: %w", err)
			}
		default:
			return new(entry.Relative), errors.New(fmt.Sprintf("Unsupported key: %s", k))
		}
	}

	return r, nil
}

var parsers = []Parser{
	&goParser{textParser: &textParser{}},
	&mdParser{textParser: &textParser{}},
	&textParser{},
}

func GetParserFor(file string) (parser Parser, err error) {
	for _, p := range parsers {
		if p.canParse(file) {
			p.SetSource(file)
			return p, nil
		}
	}

	return parser, errors.New(fmt.Sprintf("No parser found for file %s", file))
}
