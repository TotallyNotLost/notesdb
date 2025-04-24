package notesdb

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/TotallyNotLost/notesdb/parser"
)

type Notesdb struct {
	entries []*entry.Entry
}

func New() *Notesdb {
	return &Notesdb{}
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
		n.entries = append(n.entries, &entry)
	}
}

func (n *Notesdb) ListEntries() []*entry.Entry {
	return n.entries
}

func (n *Notesdb) Get(id string) (*entry.Entry, error) {
	for _, e := range n.entries {
		if e.Id == id {
			return e, nil
		}
	}

	return nil, fmt.Errorf("No entry found with id %s", id)
}
