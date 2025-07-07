package importer

import (
	"bufio"
	"os"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/TotallyNotLost/notesdb/parser"
	"github.com/TotallyNotLost/notesdb/scanner"
)

func Import(entries *map[string]*entry.Entry, source string) error {
	// TODO: Support alternative sources like URLs
	fi, err := os.Open(source)
	if err != nil {
		return err
	}
	defer fi.Close()

	r := bufio.NewReader(fi)
	scanner := scanner.New(r)
	for {
		more := scanner.Scan()
		if !more {
			break
		}

		ent, err := parser.Parse(source, scanner.Text())
		if err != nil {
			return err
		}

		for _, revision := range ent.Revisions {
			revision.Start = len(*entries)
		}

		existing, ok := (*entries)[ent.Id]

		if ok {
			// The entry has already been recorded.
			// We need to append this revision to the existing entry.
			existing.Revisions = append(existing.Revisions, ent.Revisions...)
		} else {
			(*entries)[ent.Id] = &ent
		}
	}

	return nil
}
