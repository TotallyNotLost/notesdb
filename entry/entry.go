package entry

import (
	"fmt"
	"strconv"
	"strings"
)

type Entry struct {
	Id        string     `json:"id"`
	Source    string     `json:"source"`
	Type      entryType  `json:"type"`
	Revisions []Revision `json:"revisions"`
}

type entryType int

const (
	EntryTypeText entryType = iota
	EntryTypeCode
	EntryTypeMarkdown
	EntryTypeHtml
)

type Content struct {
	Type      contentType `json:"type"`
	Value     string      `json:"value"`
	Relatives []*Relative `json:"relatives"`
}

type contentType int

const (
	ContentTypeText contentType = iota
	ContentTypeCode
	ContentTypeMarkdown
	ContentTypeLink
)

type Revision struct {
	Start     int        `json:"start"` // Inclusive
	End       int        `json:"end"`   // Exclusive
	Title     string     `json:"title"`
	Content   []Content  `json:"content"`
	Tags      []string   `json:"tags"`
	Relatives []Relative `json:"relatives"`
}

func NewRevision() Revision {
	return Revision{
		Tags:      []string{},
		Relatives: []Relative{},
	}
}

func NewEntry(id string, source string, typ entryType, title string, tags []string) Entry {
	revision := Revision{
		Title:     title,
		Tags:      tags,
		Relatives: []Relative{},
	}
	return Entry{
		Id:        id,
		Source:    source,
		Type:      typ,
		Revisions: []Revision{revision},
	}
}

func NewRelative(relationship string) (*Relative, error) {
	r := new(Relative)

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
				return new(Relative), fmt.Errorf("Relation start not an int: %w", err)
			}
		case "end":
			r.End, err = strconv.Atoi(v)
			if err != nil {
				return new(Relative), fmt.Errorf("Relation end not an int: %w", err)
			}
		default:
			return new(Relative), fmt.Errorf("Unsupported key: %s", k)
		}
	}

	return r, nil
}

type Relative struct {
	Id    string `json:"id"`
	Start int    `json:"start"` // Inclusive
	End   int    `json:"end"`   // Exclusive
}
