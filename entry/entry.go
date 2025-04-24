package entry

type Entry struct {
	Id        string      `json:"id"`
	Source    string      `json:"source"`
	Type      entryType   `json:"type"`
	Title     string      `json:"title"`
	Body      string      `json:"body"`
	Tags      []string    `json:"tags"`
	Relatives []*Relative `json:"relatives"`
}

type entryType int

const (
	EntryTypeText entryType = iota
	EntryTypeCode
	EntryTypeMarkdown
)

func NewEntry(id string, source string, typ entryType, title string, body string, tags []string) Entry {
	return Entry{
		Id:        id,
		Source:    source,
		Type:      typ,
		Title:     title,
		Body:      body,
		Tags:      tags,
		Relatives: []*Relative{},
	}
}

type Relative struct {
	Id    string `json:"id"`
	Start int    `json:"start"` // Inclusive
	End   int    `json:"end"`   // Exclusive
}
