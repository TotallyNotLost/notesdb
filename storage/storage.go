package storage

import "fmt"

type Storage struct {
	entries []*Entry
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Add(e *Entry) *Storage {
	s.entries = append(s.entries, e)
	return s
}

func (s *Storage) ListEntries() []*Entry {
	return s.entries
}

func (s *Storage) Get(id string) (*Entry, error) {
	for _, e := range s.entries {
		if e.Id == id {
			return e, nil
		}
	}

	return nil, fmt.Errorf("No entry found with id %s", id)
}
