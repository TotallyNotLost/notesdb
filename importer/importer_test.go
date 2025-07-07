package importer

import (
	"os"
	"reflect"
	"testing"

	"github.com/TotallyNotLost/notesdb/entry"
)

func TestImport(t *testing.T) {
	entries := make(map[string]*entry.Entry)
	file, _ := os.CreateTemp("", "*.md")

	defer file.Close()
	defer os.Remove(file.Name())

	data := []byte("Hello, World!\n[_metadata_:id]:# \"hello\"")
	if _, err := file.Write(data); err != nil {
		t.Error(err)
	}

	Import(&entries, file.Name())

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
				Body:      "Hello, World!\n[_metadata_:id]:# \"hello\"",
				Content:   []entry.Content{{Type: entry.ContentTypeMarkdown, Value: "Hello, World!\n[_metadata_:id]:# \"hello\""}},
				Tags:      []string{},
				Relatives: []entry.Relative{},
			},
		},
	}
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Expected\n%v,\ngot\n%v", want, *got)
	}
}
