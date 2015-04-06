package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/kurin/wordsquares/trie"
)

var matcher = regexp.MustCompile("^[a-z]*$")

var (
	size   = flag.Int("size", 2, "the word size to use")
	word   = flag.String("word", "", "the word to start with")
	dict   = flag.String("dictionary", "/usr/share/dict/words", "the dictionary to use")
	double = flag.Bool("double", false, "make a double word square")
)

func wordList(size int) ([]string, error) {
	var wl []string
	f, err := os.Open(*dict)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(scanner.Text()) == size && matcher.MatchString(scanner.Text()) {
			wl = append(wl, scanner.Text())
		}
	}
	return wl, scanner.Err()
}

func hasChar(s []byte, b byte) bool {
	for _, c := range s {
		if b == c {
			return true
		}
	}
	return false
}

func makeDoubleSquare(t *trie.Trie, words []string) [][]string {
	c := len(words)
	if c == 0 {
		return nil
	}
	s := len(words[0])
	if c == s {
		rtn := make([]string, c)
		copy(rtn, words)
		return [][]string{rtn}
	}
	matcher := make([][]byte, s)
	pfx := make([]byte, c)
	for i := 0; i < s; i++ {
		imatch := make([]byte, c)
		for j := range words {
			pfx[j] = words[j][i]
		}
		cands := t.WithPrefix(string(pfx))
		for _, cand := range cands {
			r := cand[c]
			if !hasChar(imatch, r) {
				imatch = append(imatch, r)
			}
		}
		if len(imatch) == 0 {
			return nil
		}
		matcher[i] = imatch
	}
	cWords := t.Matches(matcher)
	next := make([]string, c+1)
	copy(next, words)
	var rtn [][]string
	for _, word := range cWords {
		next[c] = word
		sqs := makeDoubleSquare(t, next)
		for _, sq := range sqs {
			rtn = append(rtn, sq)
		}
	}
	return rtn
}

func makeSquare(t *trie.Trie, words []string) [][]string {
	c := len(words)
	if c == 0 {
		return nil
	}

	s := len(words[0])
	pfx := make([]byte, c)
	for i := 0; i < s; i++ {
		for j := 0; j < c; j++ {
			pfx[j] = words[j][i]
		}
		if !t.HasPrefix(string(pfx)) {
			return nil
		}
	}

	if s == c {
		rtn := make([]string, c)
		copy(rtn, words)
		return [][]string{rtn}
	}

	for i, word := range words {
		pfx[i] = byte(word[c])
	}
	tmp := make([]string, c+1)
	copy(tmp, words)
	var all [][]string
	for _, cand := range t.WithPrefix(string(pfx)) {
		tmp[c] = cand
		sq := makeSquare(t, tmp)
		for _, s := range sq {
			all = append(all, s)
		}
	}
	return all
}

func main() {
	flag.Parse()
	squareFunc := makeSquare
	if *double {
		squareFunc = makeDoubleSquare
	}
	size := *size
	if *word != "" {
		size = len(*word)
	}
	wl, err := wordList(size)
	if err != nil {
		fmt.Println(err)
		return
	}

	t := &trie.Trie{}
	for _, w := range wl {
		t.Add(w)
	}

	if *word != "" {
		t.Add(*word)
		wSqs := squareFunc(t, []string{*word})
		for _, sq := range wSqs {
			fmt.Println("---")
			for _, a := range sq {
				fmt.Println(a)
			}
		}
		return
	}
	for _, w := range wl {
		fmt.Fprintln(os.Stderr, w)
		wSqs := squareFunc(t, []string{w})
		for _, sq := range wSqs {
			fmt.Println("---")
			for _, a := range sq {
				fmt.Println(a)
			}
		}
	}
}
