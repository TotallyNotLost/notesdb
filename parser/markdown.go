package parser

import (
	"fmt"
	"strings"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/samber/lo"
)

type mdParser struct{}

func (p mdParser) Parse(source string, text string) (e entry.Entry, err error) {
	if !strings.HasSuffix(source, ".md") {
		return e, fmt.Errorf("Yo")
	}

	revision := entry.NewRevision()
	revision.Body = text
	firstLine := lo.FirstOrEmpty(strings.Split(text, "\n"))

	if firstLine != "" {
		revision.Title = firstLine
	}

	metadata := getMetadata(text)

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

	e.Source = source
	e.Type = entry.EntryTypeMarkdown
	e.Revisions = []*entry.Revision{&revision}
	return e, nil
}
