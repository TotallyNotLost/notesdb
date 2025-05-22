package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TotallyNotLost/notesdb/entry"
)

type parser interface {
	canParse(source string) bool
	parse(source string, text string) (entry.Entry, error)
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
				return new(entry.Relative), fmt.Errorf("Relation start not an int: %w", err)
			}
		case "end":
			r.End, err = strconv.Atoi(v)
			if err != nil {
				return new(entry.Relative), fmt.Errorf("Relation end not an int: %w", err)
			}
		default:
			return new(entry.Relative), fmt.Errorf("Unsupported key: %s", k)
		}
	}

	return r, nil
}

var parsers = []parser{
	&goParser{},
	&mdParser{},
	&textParser{},
}

func Parse(source string, text string) (entry entry.Entry, err error) {
	for _, p := range parsers {
		canParse := p.canParse(source)
		if canParse {
			return p.parse(source, text)
		}
	}

	return entry, fmt.Errorf("No parser found for %s", source)
}
