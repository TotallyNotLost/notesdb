package parser

import (
	"testing"
)

func TestParseText(t *testing.T) {
	parser := &textParser{}

	t.Run("Empty string set's ID", func(t *testing.T) {
		want := "da39a3ee5e6b4b0d3255bfef95601890afd80709"
		got, _ := parser.parse("any.source.text", "")
		if got.Id != want {
			t.Errorf("want %s got %s", want, got.Id)
		}
	})
}
