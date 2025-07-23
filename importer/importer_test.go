package importer

import (
	"os"
	"reflect"
	"testing"

	"github.com/TotallyNotLost/notesdb/entry"
)

func TestImport(t *testing.T) {
	entries := make(map[string]*entry.Entry)
	dir, _ := os.MkdirTemp("", "")
	file, _ := os.CreateTemp(dir, "*.md")

	defer file.Close()
	defer os.Remove(file.Name())

	data := []byte("Hello, World!\n[_metadata_:id]:# \"hello\"")
	if _, err := file.Write(data); err != nil {
		t.Error(err)
	}

	Import(&entries, dir)

	got, ok := entries["hello"]
	if !ok {
		t.Fatal("Can't find entry for hello")
	}

	want := entry.Entry{
		Id:     "hello",
		Source: file.Name(),
		Type:   entry.EntryTypeMarkdown,
		Revisions: []entry.Revision{
			{
				Title:     "Hello, World!",
				Content:   []entry.Content{{Type: entry.ContentTypeMarkdown, Value: "Hello, World!\n[_metadata_:id]:# \"hello\""}},
				Tags:      []string{},
				Metadata:  map[string][]string{"id": {"hello"}},
				Relatives: []entry.Relative{},
			},
		},
	}
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Expected\n%v,\ngot\n%v", want, *got)
	}
}
