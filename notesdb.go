package notesdb

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/TotallyNotLost/notesdb/importer"
	"github.com/samber/lo"
)

type Notesdb struct {
	entries map[string]*entry.Entry
}

func New() *Notesdb {
	return &Notesdb{
		entries: map[string]*entry.Entry{},
	}
}

func (n *Notesdb) Import(dir string) error {
	return importer.Import(&n.entries, dir)
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

	for _, link := range revision.Metadata["link"] {
		parts := strings.Split(link, "=")
		if len(parts) < 2 || parts[0] != "id" {
			err := fmt.Errorf("%s [id=%s] Invalid link format: %s", entry.Source, entry.Id, link)
			errs = append(errs, err)
			continue
		}
		id := parts[1]

		_, ok := n.entries[id]
		if !ok {
			err := fmt.Errorf("%s [id=%s] Can't find entry with id %s", entry.Source, entry.Id, id)
			errs = append(errs, err)
		}
	}
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
		if values[i].Source != values[j].Source {
			return values[i].Source > values[j].Source
		}
		return lo.LastOrEmpty(values[i].Revisions).Start > lo.LastOrEmpty(values[j].Revisions).Start
	})

	return values
}
