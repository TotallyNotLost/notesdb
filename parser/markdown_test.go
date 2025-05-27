package parser

import (
	"testing"

	"github.com/TotallyNotLost/notesdb/entry"
	"github.com/google/go-cmp/cmp"
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
				Source: "source.md",
				Type:   2,
				Revisions: []entry.Revision{
					{
						Tags:      []string{},
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
				Id:     "the-id",
				Source: "source.md",
				Type:   2,
				Revisions: []entry.Revision{
					{
						Body: `
# Hello, World!

This is an entry with a link to [_metadata_:link]:# "id=another-id".

[_metadata_:related]:# "id=first-relative"
Testing [_metadata_:related]:# "id=second-relative"

More text

Hello [_metadata_:id]:# "the-id" world
[_metadata_:related]:# "id=third-relative"
`,
						Content: []entry.Content{
							{
								Value: `
# Hello, World!

This is an entry with a link to `},
							{Type: entry.ContentTypeLink, Value: `[_metadata_:link]:# "id=another-id"`},
							{Value: `.

[_metadata_:related]:# "id=first-relative"
Testing [_metadata_:related]:# "id=second-relative"

More text

Hello [_metadata_:id]:# "the-id" world
[_metadata_:related]:# "id=third-relative"
`},
						},
						Tags: []string{},
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
			in: `
Example entry that links to {$different-entry-id}.
			`,
			want: entry.Entry{
				Source: "source.md",
				Type:   2,
				Revisions: []entry.Revision{
					{
						Body: `
Example entry that links to [_metadata_:link]:# "id=different-entry-id".
			`,
						Tags:      []string{},
						Relatives: []entry.Relative{},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := parser.parse("source.md", tt.in)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Invalid value:\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}
