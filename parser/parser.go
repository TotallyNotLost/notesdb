package parser

import (
	"fmt"

	"github.com/TotallyNotLost/notesdb/entry"
)

type parser interface {
	canParse(source string) bool
	parse(source string, text string) (entry.Entry, error)
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
