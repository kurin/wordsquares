package trie

import "unicode/utf8"

// Trie is a trie.
type Trie struct {
	name     rune
	terminal bool
	children []*Trie
}

func (t *Trie) getNode(c rune) *Trie {
	for _, n := range t.children {
		if n.name == c {
			return n
		}
	}
	return nil
}

func (t *Trie) addChar(c rune) *Trie {
	n := t.getNode(c)
	if n == nil {
		n = &Trie{name: c}
		t.children = append(t.children, n)
	}
	return n
}

// Add adds the string to the trie.
func (t *Trie) Add(s string) {
	if len(s) == 0 {
		t.terminal = true
		return
	}
	n := t.addChar([]rune(s)[0])
	n.Add(string([]rune(s)[1:]))
}

// HasPrefix returns whether any entry in the trie has s as its prefix.
func (t *Trie) HasPrefix(s string) bool {
	if len(s) == 0 {
		return true
	}
	n := t.getNode([]rune(s)[0])
	if n == nil {
		return false
	}
	return n.HasPrefix(s[1:])
}

// HasString returns whether the given string is in the trie as a
// complete entry.
func (t *Trie) HasString(s string) bool {
	if len(s) == 0 && t.terminal {
		return true
	}
	if len(s) == 0 {
		return false
	}
	n := t.getNode([]rune(s)[0])
	if n == nil {
		return false
	}
	return n.HasString(s[1:])
}

func (t *Trie) substrings() []string {
	var strs []string
	for _, n := range t.children {
		for _, sstr := range n.substrings() {
			strs = append(strs, string(t.name)+sstr)
		}
	}
	if t.terminal {
		strs = append(strs, string(t.name))
	}
	return strs
}

func (t *Trie) subtrie(s string) *Trie {
	n := t
	for i := 0; i < utf8.RuneCountInString(s); i++ {
		n = n.getNode([]rune(s)[i])
		if n == nil {
			return nil
		}
	}
	return n
}

// WithPrefix returns all entries in the trie that begin with the
// given prefix.
func (t *Trie) WithPrefix(s string) []string {
	if len(s) == 0 {
		return t.substrings()
	}
	n := t.subtrie(s)
	var strs []string
	for _, str := range n.substrings() {
		strs = append(strs, s[:len(s)-1]+str)
	}
	return strs
}
