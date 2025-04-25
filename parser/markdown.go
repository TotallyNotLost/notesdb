package parser

import (
	"fmt"
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

	revision := lo.FirstOrEmpty(e.Revisions)
	firstLine := lo.FirstOrEmpty(strings.Split(revision.Body, "\n"))

	if firstLine != "" {
		revision.Title = firstLine
	}

	metadata := getMetadata(revision.Body)

	if ids, ok := metadata["id"]; ok {
		e.Id = lo.LastOrEmpty(ids)
	}
	if relatives, ok := metadata["related"]; ok {
		for _, relationship := range relatives {
			relative, err := NewRelative(relationship)
			if err != nil {
				err = fmt.Errorf("[id=%s] Processing related metadata: %w", e.Id, err)
				return entry.Entry{}, err
			}
			revision.Relatives = append(revision.Relatives, relative)
			if strings.HasPrefix(relative.Id, "#") {
				revision.Tags = append(revision.Tags, strings.TrimPrefix(relative.Id, "#"))
			}
		}

	}

	e.Type = entry.EntryTypeMarkdown
	return e, nil
}
