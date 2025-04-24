package parser

import (
	"io"
	"strings"

	"github.com/TotallyNotLost/notesdb/storage"
	"github.com/samber/lo"
)

type mdParser struct {
	textParser *textParser
}

func (p *mdParser) SetSource(source string)     { p.textParser.SetSource(source) }
func (p *mdParser) SetReader(reader io.Reader)  { p.textParser.SetReader(reader) }
func (p *mdParser) canParse(source string) bool { return strings.HasSuffix(source, ".md") }
func (p *mdParser) Next() (entry storage.Entry, err error) {
	entry, err = p.textParser.Next()
	if err != nil {
		return entry, err
	}

	firstLine := lo.FirstOrEmpty(strings.Split(entry.Body, "\n"))

	if firstLine != "" {
		entry.Title = firstLine
	}

	metadata := getMetadata(entry.Body)

	if ids, ok := metadata["id"]; ok {
		entry.Id = lo.LastOrEmpty(ids)
	}
	if relatives, ok := metadata["related"]; ok {
		for _, relationship := range relatives {
			relative, err := NewRelative(relationship)
			if err != nil {
				return storage.Entry{}, err
			}
			entry.Relatives = append(entry.Relatives, relative)
			if strings.HasPrefix(relative.Id, "#") {
				entry.Tags = append(entry.Tags, strings.TrimPrefix(relative.Id, "#"))
			}
		}

	}

	entry.Type = storage.EntryTypeMarkdown
	return entry, nil
}
