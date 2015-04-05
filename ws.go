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
	size = flag.Int("size", 2, "the word size to use")
	word = flag.String("word", "", "the word to start with")
	dict = flag.String("dictionary", "/usr/share/dict/words", "the dictionary to use")
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

func makeSquare(words []string, t *trie.Trie) [][]string {
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
		sq := makeSquare(tmp, t)
		for _, s := range sq {
			all = append(all, s)
		}
	}
	return all
}

func main() {
	flag.Parse()
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
		wSqs := makeSquare([]string{*word}, t)
		for _, sq := range wSqs {
			fmt.Println("---")
			for _, a := range sq {
				fmt.Println(a)
			}
		}
		return
	}
	for _, w := range wl {
		wSqs := makeSquare([]string{w}, t)
		for _, sq := range wSqs {
			fmt.Println("---")
			for _, a := range sq {
				fmt.Println(a)
			}
		}
	}
}
