package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/kurin/wordsquares/trie"
)

var matcher = regexp.MustCompile("^[a-z]*$")

func wordList(size int) ([]string, error) {
	var wl []string
	f, err := os.Open("/usr/share/dict/words")
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

func makeSquare(words []string, t *trie.Trie) []string {
	c := len(words)
	if c == 0 {
		return nil
	}

	s := len(words[0])
	b := make([]byte, c)
	for i := 0; i < s; i++ {
		for j := 0; j < c; j++ {
			b[j] = words[j][i]
		}
		if !t.HasPrefix(string(b)) {
			return nil
		}
	}

	if s == c {
		return words
	}

	pfx := make([]byte, c)
	for i, word := range words {
		pfx[i] = byte(word[c])
	}
	tmp := make([]string, len(words)+1)
	copy(tmp, words)
	for _, cand := range t.WithPrefix(string(pfx)) {
		tmp[len(words)] = cand
		if sq := makeSquare(tmp, t); sq != nil {
			return sq
		}
	}
	return nil
}

func main() {
	wl, err := wordList(8)
	if err != nil {
		fmt.Println(err)
		return
	}

	t := &trie.Trie{}
	for _, w := range wl {
		t.Add(w)
	}
	for _, w := range wl {
		fmt.Println(w)
		if ans := makeSquare([]string{w}, t); ans != nil {
			for _, a := range ans {
				fmt.Println(a)
			}
			return
		}
	}
}
