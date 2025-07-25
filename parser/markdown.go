package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/samber/lo"
)

type mdParser struct{}

func (p mdParser) canParse(source string) bool { return strings.HasSuffix(source, ".md") }
func (p mdParser) parse(text string, e *entry.Entry) error {
	text = expandShortLinks(text)
	revision := entry.NewRevision()
	revision.Content = append(revision.Content, getContent(text)...)
	firstLine := lo.FirstOrEmpty(strings.Split(text, "\n"))

	if firstLine != "" {
		revision.Title = firstLine
	}

	revision.Metadata = getMetadata(text)

	if ids, ok := revision.Metadata["id"]; ok {
		e.Id = lo.LastOrEmpty(ids)
	}
	if relatives, ok := revision.Metadata["related"]; ok {
		for _, relationship := range relatives {
			relative, err := entry.NewRelative(relationship)
			if err != nil {
				continue
			}
			revision.Relatives = append(revision.Relatives, *relative)
			if strings.HasPrefix(relative.Id, "#") {
				revision.Tags = append(revision.Tags, strings.TrimPrefix(relative.Id, "#"))
			}
		}

	}
	// TODO: Relatives should include entries linked to within the text.

	e.Type = entry.EntryTypeMarkdown
	e.Revisions = []entry.Revision{revision}
	return nil
}

func getContent(text string) []entry.Content {
	if text == "" {
		return []entry.Content{}
	}

	contents := []entry.Content{}
	r := regexp.MustCompile("\\[_metadata_:link]:# \"([^\"]*)\"")

	text = r.ReplaceAllStringFunc(text, func(match string) string {
		submatches := r.FindStringSubmatch(match)
		return fmt.Sprintf("[<notesdb>](%s)", submatches[1])
	})

	contents = append(contents, entry.Content{Type: entry.ContentTypeMarkdown, Value: text})

	return contents
}

func expandShortLinks(text string) string {
	r, _ := regexp.Compile("\\{\\$([^}]*)\\}")

	return r.ReplaceAllStringFunc(text, func(match string) string {
		id := r.FindStringSubmatch(match)[1]

		return fmt.Sprintf("[_metadata_:link]:# \"id=%s\"", id)
	})
}

func getMetadata(text string) map[string][]string {
	lines := strings.Split(text, "\n")

	o := make(map[string][]string)

	for _, l := range lines {
		r := regexp.MustCompile("\\[_metadata_:*(\\w+)\\]:# \"([^\"]*)\"")

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
