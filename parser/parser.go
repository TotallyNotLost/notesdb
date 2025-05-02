package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/TotallyNotLost/notesdb/entry"
)

type Parser interface {
	Parse(source string, text string) (entry.Entry, error)
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

var parsers = []Parser{
	&goParser{},
	&mdParser{},
	&textParser{},
}

func Parse(source string, text string) (entry entry.Entry, err error) {
	for _, parser := range parsers {
		entry, err = parser.Parse(source, text)
		if err == nil {
			return entry, nil
		}
	}

	return entry, fmt.Errorf("No parser found for %s", source)
}
