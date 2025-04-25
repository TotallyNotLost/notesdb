package notesdb

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"iter"
	"maps"
	"os"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/TotallyNotLost/notesdb/parser"
)

type Notesdb struct {
	entries map[string]*entry.Entry
}

func New() *Notesdb {
	return &Notesdb{
		entries: make(map[string]*entry.Entry),
	}
}

func (n *Notesdb) Import(source string) {
	// TODO: Support alternative sources like URLs
	fi, err := os.Open(source)
	if err != nil {
		// TODO: Don't use panic
		panic(err)
	}
	defer fi.Close()

	p, err := parser.GetParserFor(source)
	if err != nil {
		// TODO: Don't use panic
		panic(err)
	}
	p.SetReader(bufio.NewReader(fi))
	for {
		entry, err := p.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			// TODO: Don't use panic
			panic(err)
		}

		existing, ok := n.entries[entry.Id]

		if ok {
			// The entry has already been recorded.
			// We need to append this revision to the existing entry.
			existing.Revisions = append(existing.Revisions, entry.Revisions...)
		} else {
			n.entries[entry.Id] = &entry
		}
	}
}

/*
 * Verify that the imports were all successful.
 *
 * Returns an error if any references couldn't be resolved.
 */
func (n *Notesdb) Verify() error {
	var errs []error

	for entry := range n.All() {
		err := n.verifyEntry(entry)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (n *Notesdb) verifyEntry(entry *entry.Entry) error {
	var errs []error

	for _, revision := range entry.Revisions {
		err := n.verifyRevision(entry, revision)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (n *Notesdb) verifyRevision(entry *entry.Entry, revision *entry.Revision) error {
	var errs []error

	for _, relative := range revision.Relatives {
		_, ok := n.entries[relative.Id]
		if !ok {
			err := fmt.Errorf("Can't find entry with id %s referenced by %s in %s", relative.Id, entry.Id, entry.Source)
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (n *Notesdb) All() iter.Seq[*entry.Entry] {
	return maps.Values(n.entries)
}
