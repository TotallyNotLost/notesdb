package parser

import (
	"fmt"
	"regexp"
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
