// Deals with Burrows-Wheeler transformations and indexing.
package bwt

import (
	"bytes"
	"fmt"
	"sort"
)

// Returns a BW-transform of the given sequence, without changing it. Adds a $
// sign as a terminating character.
func Transform(str []byte) []byte {
	return newBwSorter(str).sort().transform()
}

// A sorting interface for a sequence's suffixes.
type bwSorter struct {
	str  []byte  // Sequence to transform.
	perm []int   // Suffix permutation to sort.
}

// Creates a new sorter for the given sequence.
func newBwSorter(str []byte) *bwSorter {
	// Copy str and add a terminator character.
	strCopy := make([]byte, len(str) + 1)
	copy(strCopy, str)
	strCopy[len(strCopy) - 1] = '$'
	
	perm := make([]int, len(strCopy))
	for i := range perm {
		perm[i] = i
	}
	return &bwSorter{strCopy, perm}
}

// Sorts the sorter's permutation.
func (b *bwSorter) sort() *bwSorter {
	sort.Sort(b)
	return b
}

// Sorting functions.
func (b *bwSorter) Len() int {
	return len(b.perm)
}
func (b *bwSorter) Less(i, j int) bool {
	return bytes.Compare(b.str[b.perm[i]:], b.str[b.perm[j]:]) < 0
}
func (b *bwSorter) Swap(i, j int) {
	b.perm[i], b.perm[j] = b.perm[j], b.perm[i]
}

// Applies BW-transform of the sorter's sequence. Should be called after
// sorting. Does not change the input sequence.
func (b *bwSorter) transform() []byte {
	arr := make([]byte, len(b.str))
	for i := range arr {
		arr[i] = b.str[(b.perm[i]-1+len(b.str))%len(b.str)]
	}
	return arr
}

// Indexes the ranks of characters in a BW-transformed sequence.
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
