package parser

import (
	"io"
	"strings"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/samber/lo"
)

type mdParser struct {
	textParser *textParser
}

func (p *mdParser) SetSource(source string)     { p.textParser.SetSource(source) }
func (p *mdParser) SetReader(reader io.Reader)  { p.textParser.SetReader(reader) }
func (p *mdParser) canParse(source string) bool { return strings.HasSuffix(source, ".md") }
func (p *mdParser) Next() (e entry.Entry, err error) {
	e, err = p.textParser.Next()
	if err != nil {
		return e, err
	}

	firstLine := lo.FirstOrEmpty(strings.Split(e.Body, "\n"))

	if firstLine != "" {
		e.Title = firstLine
	}

	metadata := getMetadata(e.Body)

	if ids, ok := metadata["id"]; ok {
		e.Id = lo.LastOrEmpty(ids)
	}
	if relatives, ok := metadata["related"]; ok {
		for _, relationship := range relatives {
			relative, err := NewRelative(relationship)
			if err != nil {
				return entry.Entry{}, err
			}
			e.Relatives = append(e.Relatives, relative)
			if strings.HasPrefix(relative.Id, "#") {
				e.Tags = append(e.Tags, strings.TrimPrefix(relative.Id, "#"))
			}
		}

	}

	e.Type = entry.EntryTypeMarkdown
	return e, nil
}
