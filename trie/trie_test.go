package trie

import (
	"reflect"
	"sort"
	"testing"
)

func TestHasPrefix(t *testing.T) {
	table := []struct {
		dict  []string
		tests map[string]bool
	}{
		{
			dict: []string{"this", "is", "a", "test"},
			tests: map[string]bool{
				"th":  true,
				"thi": true,
				"ok":  false,
			},
		},
	}
	for _, e := range table {
		tr := &Trie{}
		for _, w := range e.dict {
			tr.Add(w)
		}
		for pfx, want := range e.tests {
			got := tr.HasPrefix(pfx)
			if got != want {
				t.Errorf("tr.HasPrefix(%q): got %v, want %v", pfx, got, want)
			}
		}
	}
}

func TestHasString(t *testing.T) {
	table := []struct {
		dict  []string
		tests map[string]bool
	}{
		{
			dict: []string{"this", "is", "a", "test", "is"},
			tests: map[string]bool{
				"th":  false,
				"thi": false,
				"ok":  false,
				"is":  true,
			},
		},
	}
	for _, e := range table {
		tr := &Trie{}
		for _, w := range e.dict {
			tr.Add(w)
		}
		for str, want := range e.tests {
			got := tr.HasString(str)
			if got != want {
				t.Errorf("tr.HasString(%q): got %v, want %v", str, got, want)
			}
		}
	}
}

func TestWithPrefix(t *testing.T) {
	table := []struct {
		dict  []string
		tests map[string][]string
	}{
		{
			dict: []string{"this", "that", "other", "tank", "think"},
			tests: map[string][]string{
				"th": []string{"this", "that", "think"},
			},
		},
	}
	for _, e := range table {
		tr := &Trie{}
		for _, w := range e.dict {
			tr.Add(w)
		}
		for str, want := range e.tests {
			got := tr.WithPrefix(str)
			sort.Strings(got)
			sort.Strings(want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("tr.WithPrefix(%q): got %v, want %v", str, got, want)
			}
		}
	}
}
