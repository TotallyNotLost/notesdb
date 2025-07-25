package notesdb

import (
	"os"
	"sort"
	"testing"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
)

func TestImport_importsAllEntries(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []string
	}{
		"Empty input": {
			input: "",
			want:  []string{},
		},
		"Single entry": {
			input: `First entry [_metadata_:id]:# "first"`,
			want:  []string{"first"},
		},
		"Multiple entries": {
			input: "First entry [_metadata_:id]:# \"first\"\n---\nSecond entry [_metadata_:id]:# \"second\"",
			want:  []string{"first", "second"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db := New()

			dir, _ := os.MkdirTemp("", "")
			file, _ := os.CreateTemp(dir, "*.md")

			defer file.Close()
			defer os.Remove(file.Name())
			defer os.Remove(dir)

			data := []byte(tt.input)
			if _, err := file.Write(data); err != nil {
				t.Error(err)
			}

			db.Import(dir)

			got := lo.Map(lo.Values(db.entries), func(e *entry.Entry, index int) string {
				return e.Id
			})
			sort.Strings(got)

			if !cmp.Equal(tt.want, got) {
				diff := cmp.Diff(tt.want, got)
				t.Errorf("\nexpected:\n%v\n\ngot:\n%v\n\ndiff:\n%v", tt.want, got, diff)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	tests := map[string]struct {
		input     string
		wantError bool
	}{
		"Empty input": {
			input:     "",
			wantError: false,
		},
		"Valid link": {
			input:     "First {$second}\n[_metadata_:id]:# \"first\"\n---\nSecond\n[_metadata_:id]:# \"second\"",
			wantError: false,
		},
		"Link to non-existent entry": {
			input:     "First {$third}\n[_metadata_:id]:# \"first\"\n---\nSecond\n[_metadata_:id]:# \"second\"",
			wantError: true,
		},
		"Valid related": {
			input:     "First\n[_metadata_:id]:# \"first\"\n[_metadata_:related]:# \"id=second\"\n---\nSecond\n[_metadata_:id]:# \"second\"",
			wantError: false,
		},
		"Related to non-existent entry": {
			input:     "First\n[_metadata_:id]:# \"first\"\n[_metadata_:related]:# \"id=third\"\n---\nSecond\n[_metadata_:id]:# \"second\"",
			wantError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db := New()

			dir, _ := os.MkdirTemp("", "")
			file, _ := os.CreateTemp(dir, "*.md")

			defer file.Close()
			defer os.Remove(file.Name())
			defer os.Remove(dir)

			data := []byte(tt.input)
			if _, err := file.Write(data); err != nil {
				t.Error(err)
			}

			db.Import(dir)

			got := db.Verify()
			gotError := got != nil

			if gotError != tt.wantError {
				t.Errorf("expected error?: %v, got error?: %v. Error was %v", tt.wantError, gotError, got)
			}
		})
	}
}
