package trie

// Trie is a trie.  It doesn't support unicode strings.
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

// Add adds the string to the trie.
func (t *Trie) Add(s string) {
	if len(s) == 0 {
		t.terminal = true
		return
	}
	n := t.addChar(s[0])
	n.Add(s[1:])
}

// HasPrefix returns whether any entry in the trie has s as its prefix.
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

// HasString returns whether the given string is in the trie as a
// complete entry.
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

func (t *Trie) subbytes(soFar []byte, strs *[]string) {
	soFar = append(soFar, t.name)
	for _, n := range t.children {
		n.subbytes(soFar, strs)
	}
	if t.terminal {
		*strs = append(*strs, string(soFar))
	}
}

func (t *Trie) subtrie(s string) *Trie {
	n := t
	for i := 0; i < len(s); i++ {
		n = n.getNode(s[i])
		if n == nil {
			return nil
		}
	}
	return n
}

// WithPrefix returns all entries in the trie that begin with the
// given prefix.
func (t *Trie) WithPrefix(s string) []string {
	var strs []string
	if len(s) == 0 {
		t.subbytes([]byte{}, &strs)
		return strs
	}
	n := t.subtrie(s)
	if n == nil {
		return nil
	}
	for _, n := range n.children {
		n.subbytes([]byte(s), &strs)
	}
	return strs
}

// Previous WithPrefix, kept for benchmark purposes.
func (t *Trie) oldWithPrefix(s string) []string {
	if len(s) == 0 {
		return t.substrings()
	}
	n := t.subtrie(s)
	if n == nil {
		return nil
	}
	var strs []string
	for _, str := range n.substrings() {
		strs = append(strs, s[:len(s)-1]+str)
	}
	return strs
}

// Matches returns all the entries in the trie which match the given byte list.
// All the returned strings will have one character from bl[0] in the first
// position, bl[1] in the second position, etc.
func (t *Trie) Matches(bl [][]byte) []string {
	var name string
	if t.name > 0 {
		// the 0 code point is ""
		name = string(t.name)
	}
	if len(bl) == 0 && t.terminal {
		return []string{name}
	}
	if len(bl) == 0 {
		return nil
	}
	var rtn []string
	for _, r := range bl[0] {
		n := t.getNode(r)
		if n == nil {
			continue
		}
		for _, mat := range n.Matches(bl[1:]) {
			rtn = append(rtn, name+mat)
		}
	}
	return rtn
}
