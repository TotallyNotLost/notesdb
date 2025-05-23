package notesdb

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/TotallyNotLost/notesdb/parser"
	"github.com/TotallyNotLost/notesdb/scanner"
	"github.com/samber/lo"
)

type Notesdb struct {
	entries map[string]*entry.Entry
}

func New() *Notesdb {
	return &Notesdb{
		entries: make(map[string]*entry.Entry),
	}
}

func (n *Notesdb) Import(source string) error {
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
			revision.Start = len(n.entries)
		}

		existing, ok := n.entries[ent.Id]

		if ok {
			// The entry has already been recorded.
			// We need to append this revision to the existing entry.
			existing.Revisions = append(existing.Revisions, ent.Revisions...)
		} else {
			n.entries[ent.Id] = &ent
		}
	}

	return nil
}

/*
 * Verify that the imports were all successful.
 *
 * Returns an error if any references couldn't be resolved.
 */
func (n *Notesdb) Verify() error {
	var errs []error

	for _, entry := range n.All() {
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
		err := n.verifyRevision(entry, &revision)
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
			err := fmt.Errorf("%s [id=%s] Can't find entry with id %s", entry.Source, entry.Id, relative.Id)
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (n *Notesdb) All() []*entry.Entry {
	var values []*entry.Entry
	for _, value := range n.entries {
		values = append(values, value)
	}

	sort.Slice(values, func(i, j int) bool {
		return lo.LastOrEmpty(values[i].Revisions).Start > lo.LastOrEmpty(values[j].Revisions).Start
	})

	return values
}
