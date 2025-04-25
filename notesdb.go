package notesdb

import (
	"bufio"
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

func (n *Notesdb) All() iter.Seq[*entry.Entry] {
	return maps.Values(n.entries)
}

func (n *Notesdb) Get(id string) (*entry.Entry, error) {
	for _, e := range n.entries {
		if e.Id == id {
			return e, nil
		}
	}

	return nil, fmt.Errorf("No entry found with id %s", id)
}
