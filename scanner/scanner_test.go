package scanner

import (
	"slices"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := map[string]struct {
		in   string
		want []string
	}{
		"Empty string": {
			in:   "",
			want: []string{},
		},
		"No newlines": {
			in:   "Testing",
			want: []string{"Testing"},
		},
		"Newlines, but no separators": {
			in:   "Testing\n123\n",
			want: []string{"Testing\n123\n"},
		},
		"Incomplete separator (\n---)": {
			in:   "Testing\n---123",
			want: []string{"Testing\n---123"},
		},
		"Incomplete separator (---\n)": {
			in:   "Testing---\n123",
			want: []string{"Testing---\n123"},
		},
		"Complete separator": {
			in:   "Testing\n---\n123",
			want: []string{"Testing", "123"},
		},
		"Multiple separators": {
			in:   "Testing\n---\n123\n---\nMore testing",
			want: []string{"Testing", "123", "More testing"},
		},
		"Newlines surrounding complete separator": {
			in:   "Testing\n\n---\n\n123",
			want: []string{"Testing\n", "\n123"},
		},
		"Only a separator": {
			in:   "\n---\n",
			want: []string{""},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got []string
			r := strings.NewReader(test.in)
			s := New(r)

			for s.Scan() {
				got = append(got, s.Text())
			}

			if !slices.Equal(got, test.want) {
				t.Errorf("Expected %s, got %s", test.want, got)
			}
		})
	}
}
