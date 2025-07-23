package parser

import "testing"

func TestParse_setsId(t *testing.T) {
	tests := map[string]struct {
		source string
		text   string
		want   string
	}{
		"Markdown file with empty contents": {source: "example.md", text: "", want: "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		"Markdown file with no metadata":    {source: "example.md", text: "Hello, World!", want: "0a0a9f2a6772942557ab5355d76af442f8f65e01"},
		"Markdown file with metadata":       {source: "example.md", text: `[_metadata_:link]:# "test"`, want: "6a494bff2f4230b5bc5ea07adcb257ba5f677474"},
		"Markdown file with id in metadata": {source: "example.md", text: `[_metadata_:id]:# "test"`, want: "test"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := Parse(tt.source, tt.text)
			if got.Id != tt.want {
				t.Fatalf("Expected %s, got %s", tt.want, got.Id)
			}
		})
	}
}
