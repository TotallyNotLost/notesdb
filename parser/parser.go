package parser

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/TotallyNotLost/notesdb/entry"
)

type parser interface {
	canParse(source string) bool
	parse(text string, e *entry.Entry) error
}

var parsers = []parser{
	&goParser{},
	&mdParser{},
	&textParser{},
}

func Parse(source string, text string) (entry entry.Entry, err error) {
	entry.Source = source
	h := sha1.New()
	h.Write([]byte(text))
	entry.Id = hex.EncodeToString(h.Sum(nil))
	for _, p := range parsers {
		canParse := p.canParse(source)
		if canParse {
			err = p.parse(text, &entry)
			return
		}
	}

	err = fmt.Errorf("No parser found for %s", source)
	return
}
