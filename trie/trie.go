package trie

// Trie is a radix tree.
type Trie struct {
	name     byte
	terminal bool
	children []*Trie
}

func (t *Trie) getNode(c byte) *Trie {
	for _, n := range t.children {
		if n.name == c {
			return n
		}
	}
	return nil
}

func (t *Trie) addChar(c byte) *Trie {
	n := t.getNode(c)
	if n == nil {
		n = &Trie{name: c}
		t.children = append(t.children, n)
	}
	return n
}

func (t *Trie) Add(s string) {
	if len(s) == 0 {
		t.terminal = true
		return
	}
	n := t.addChar(s[0])
	n.Add(s[1:])
}

func (t *Trie) HasPrefix(s string) bool {
	if len(s) == 0 {
		return true
	}
	n := t.getNode(s[0])
	if n == nil {
		return false
	}
	return n.HasPrefix(s[1:])
}

func (t *Trie) HasString(s string) bool {
	if len(s) == 0 && t.terminal {
		return true
	}
	if len(s) == 0 {
		return false
	}
	n := t.getNode(s[0])
	if n == nil {
		return false
	}
	return n.HasString(s[1:])
}

func (t *Trie) Substrings() []string {
	var strs []string
	for _, n := range t.children {
		for _, sstr := range n.Substrings() {
			strs = append(strs, string(t.name)+sstr)
		}
	}
	if t.terminal {
		strs = append(strs, string(t.name))
	}
	return strs
}

func (t *Trie) Subtrie(s string) *Trie {
	n := t
	for i := 0; i < len(s); i++ {
		n = n.getNode(s[i])
		if n == nil {
			return nil
		}
	}
	return n
}

func (t *Trie) WithPrefix(s string) []string {
	n := t.Subtrie(s)
	var strs []string
	for _, str := range n.Substrings() {
		strs = append(strs, s[:len(s)-1]+str)
	}
	return strs
}
