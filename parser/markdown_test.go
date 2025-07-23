package parser

import (
	"testing"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
)

func TestCanParseMarkdown(t *testing.T) {
	parser := &mdParser{}

	tests := map[string]struct {
		in   string
		want bool
	}{
		"Empty string": {
			in:   "",
			want: false,
		},
		"No extension": {
			in:   "README",
			want: false,
		},
		"Unsupported extension": {
			in:   "README.txt",
			want: false,
		},
		".md extension": {
			in:   "README.md",
			want: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := parser.canParse(tt.in)
			if got != tt.want {
				t.Errorf("Expected %t got %t", tt.want, got)
			}
		})
	}
}

func TestParseMarkdown(t *testing.T) {
	parser := &mdParser{}

	tests := map[string]struct {
		in   string
		want entry.Entry
	}{
		"Empty string": {
			in: "",
			want: entry.Entry{
				Type: entry.EntryTypeMarkdown,
				Revisions: []entry.Revision{
					{
						Tags:      []string{},
						Metadata:  map[string][]string{},
						Relatives: []entry.Relative{},
					},
				},
			},
		},
		"With metadata": {
			in: `
# Hello, World!

This is an entry with a link to [_metadata_:link]:# "id=another-id".

[_metadata_:related]:# "id=first-relative"
Testing [_metadata_:related]:# "id=second-relative"

More text

Hello [_metadata_:id]:# "the-id" world
[_metadata_:related]:# "id=third-relative"
`,
			want: entry.Entry{
				Id:   "the-id",
				Type: entry.EntryTypeMarkdown,
				Revisions: []entry.Revision{
					{
						Content: []entry.Content{
							{
								Type: entry.ContentTypeMarkdown,
								Value: `
# Hello, World!

This is an entry with a link to `},
							{Type: entry.ContentTypeLink, Value: `[_metadata_:link]:# "id=another-id"`},
							{Type: entry.ContentTypeMarkdown, Value: `.

[_metadata_:related]:# "id=first-relative"
Testing [_metadata_:related]:# "id=second-relative"

More text

Hello [_metadata_:id]:# "the-id" world
[_metadata_:related]:# "id=third-relative"
`},
						},
						Tags: []string{},
						Metadata: map[string][]string{
							"id":      {"the-id"},
							"link":    {"id=another-id"},
							"related": {"id=first-relative", "id=second-relative", "id=third-relative"},
						},
						Relatives: []entry.Relative{
							{
								Id: "first-relative",
							},
							{
								Id: "second-relative",
							},
							{
								Id: "third-relative",
							},
						},
					},
				},
			},
		},
		"With short links": {
			in: `Example entry that links to {$different-entry-id}.`,
			want: entry.Entry{
				Type: entry.EntryTypeMarkdown,
				Revisions: []entry.Revision{
					{
						Title: `Example entry that links to [_metadata_:link]:# "id=different-entry-id".`,
						Content: []entry.Content{
							{Type: entry.ContentTypeMarkdown, Value: "Example entry that links to "},
							{Type: entry.ContentTypeLink, Value: `[_metadata_:link]:# "id=different-entry-id"`},
							{Type: entry.ContentTypeMarkdown, Value: "."},
						},
						Metadata:  map[string][]string{"link": {"id=different-entry-id"}},
						Tags:      []string{},
						Relatives: []entry.Relative{},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var e entry.Entry
			parser.parse(tt.in, &e)
			if !cmp.Equal(e, tt.want) {
				t.Errorf("Invalid value:\n%s", cmp.Diff(tt.want, e))
			}
		})
	}
}

func TestParseMarkdown_setsMetadata(t *testing.T) {
	parser := &mdParser{}

	tests := map[string]struct {
		in   string
		want map[string][]string
	}{
		"Empty string": {
			in:   "",
			want: map[string][]string{},
		},
		"No metadata": {
			in:   "Hello, World!",
			want: map[string][]string{},
		},
		"With metadata": {
			in: `
# Hello, World!

This is an entry with a link to [_metadata_:link]:# "id=another-id".

[_metadata_:related]:# "id=first-relative"
Testing [_metadata_:related]:# "id=second-relative"

More text

Hello [_metadata_:id]:# "the-id" world
[_metadata_:related]:# "id=third-relative"
`,
			want: map[string][]string{
				"id":      {"the-id"},
				"link":    {"id=another-id"},
				"related": {"id=first-relative", "id=second-relative", "id=third-relative"},
			},
		},
		"With short links": {
			in:   `Example entry that links to {$different-entry-id}.`,
			want: map[string][]string{"link": {"id=different-entry-id"}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var e entry.Entry
			parser.parse(tt.in, &e)
			revision := lo.LastOrEmpty(e.Revisions)
			if !cmp.Equal(revision.Metadata, tt.want) {
				t.Errorf("Invalid value:\n%s", cmp.Diff(tt.want, revision.Metadata))
			}
		})
	}
}
