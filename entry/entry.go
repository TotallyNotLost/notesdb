package entry

type Entry struct {
	Id        string      `json:"id"`
	Source    string      `json:"source"`
	Type      entryType   `json:"type"`
	Revisions []*Revision `json:"revisions"`
}

type entryType int

const (
	EntryTypeText entryType = iota
	EntryTypeCode
	EntryTypeMarkdown
)

type Revision struct {
	Title     string      `json:"title"`
	Body      string      `json:"body"`
	Tags      []string    `json:"tags"`
	Relatives []*Relative `json:"relatives"`
}

func NewEntry(id string, source string, typ entryType, title string, body string, tags []string) Entry {
	revision := &Revision{
		Title:     title,
		Body:      body,
		Tags:      tags,
		Relatives: []*Relative{},
	}
	return Entry{
		Id:        id,
		Source:    source,
		Type:      typ,
		Revisions: []*Revision{revision},
	}
}

type Relative struct {
	Id    string `json:"id"`
	Start int    `json:"start"` // Inclusive
	End   int    `json:"end"`   // Exclusive
}
