package main

import (
	"bytes"
	"fmt"
	"sort"
)

func main() {
	s := "how much wood would a wood chuck chuck"
	fmt.Println("'" + string(transform([]byte(s))) + "'")
	fmt.Println(newRankIndex([]byte("AATACGCTGA"), 3))
	fmt.Println("AATACGCTGA")
}

func transform(str []byte) []byte {
	return newBwSorter(str).sort().transform()
}

type bwSorter struct {
	str  []byte
	perm []int
}

func newBwSorter(str []byte) *bwSorter {
	perm := make([]int, len(str))
	for i := range perm {
		perm[i] = i
	}
	return &bwSorter{str, perm}
}

func (b *bwSorter) sort() *bwSorter {
	sort.Sort(b)
	return b
}

func (b *bwSorter) Len() int {
	return len(b.perm)
}
func (b *bwSorter) Less(i, j int) bool {
	return bytes.Compare(b.str[b.perm[i]:], b.str[b.perm[j]:]) < 0
}
func (b *bwSorter) Swap(i, j int) {
	b.perm[i], b.perm[j] = b.perm[j], b.perm[i]
}

func (b *bwSorter) transform() []byte {
	arr := make([]byte, len(b.str))
	for i := range arr {
		arr[i] = b.str[(b.perm[i]-1+len(b.str))%len(b.str)]
	}
	return arr
}

type rankIndex struct {
	chars map[byte]int // From char to index in rank.
	ranks [][]int      // Ranks[i] is ranks at i*jump including.
	jump  int          // How often to save ranks.
}

func newRankIndex(str []byte, jump int) *rankIndex {
	result := &rankIndex{nil, nil, jump}

	// Create a map of characters to index in the rank arrays.
	chars := make(map[byte]int)
	for _, b := range str {
		if chars[b] == 0 {
			chars[b] = len(chars) + 1
		}
	}

	// Indexes were +1 just to avoid zeros.
	for b := range chars {
		chars[b]--
	}

	result.chars = chars

	// Create ranks array (+-1 because the first out of each jump gets
	// an array).
	ranks := make([][]int, (len(str)-1)/jump+1)

	// Count ranks in sequence.
	rank := make([]int, len(chars))
	for i, b := range str {
		// Increase rank.
		rank[chars[b]]++

		// Copy array if hit a checkpoint.
		if i%jump == 0 {
			cp := make([]int, len(rank))
			copy(cp, rank)
			ranks[i/jump] = cp
		}
	}

	result.ranks = ranks

	return result
}

func (r *rankIndex) String() string {
	buf := bytes.NewBuffer(nil)

	fmt.Fprint(buf, "chars: ")
	for b, i := range r.chars {
		fmt.Fprintf(buf, "%d-%c ", i, b)
	}
	fmt.Fprintln(buf)

	fmt.Fprintln(buf, "jump:", r.jump)

	fmt.Fprintln(buf, "ranks:")
	for _, rank := range r.ranks {
		for _, r := range rank {
			fmt.Fprintf(buf, "%d\t", r)
		}
		fmt.Fprintln(buf)
	}

	return buf.String()
}

func (r *rankIndex) rankOf(char byte, at int, str []byte) int {
	// Ranks at the closest rank array.
	result := r.ranks[at/r.jump][r.chars[char]]

	// Count up to 'at'.
	for i := at/r.jump*r.jump + 1; i <= at; i++ {
		if str[i] == char {
			result++
		}
	}

	return result
}
