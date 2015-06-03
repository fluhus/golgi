// Deals with Burrows-Wheeler transformations and indexing.
package bwt

import (
	"bytes"
	"sort"
)


// ----- TRANSFORM -------------------------------------------------------------

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

// Sorting interface functions.
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


// ----- RANK INDEX ------------------------------------------------------------

// Indexes the ranks of characters in a BW-transformed sequence.
// Ranks are 1-base, since they indicate counts of character appearances.
type rankIndex struct {
	str   []byte       // Input sequence.
	chars map[byte]int // From char to index in rank.
	ranks [][]int      // Ranks[i] is ranks at i*jump including.
	jump  int          // How often to save ranks.
}

// Creates an index for the given BW-transformed sequence.
func newRankIndex(str []byte, jump int) *rankIndex {
	result := &rankIndex{str, nil, nil, jump}

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
	ranks := make([][]int, (len(str)-1) / jump+1)

	// Count ranks in sequence.
	rank := make([]int, len(chars))
	for i, b := range str {
		// Increase rank.
		rank[chars[b]]++

		// Copy array if hit a checkpoint.
		if i % jump == 0 {
			cp := make([]int, len(rank))
			copy(cp, rank)
			ranks[i / jump] = cp
		}
	}

	result.ranks = ranks

	return result
}

// Returns the rank of a character at a given position, including
// that position.
func (r *rankIndex) rankOf(char byte, at int) int {
	// Check if character exists.
	if _, ok := r.chars[char]; !ok {
		return 0
	}

	// Ranks at the closest rank array.
	result := r.ranks[at / r.jump][r.chars[char]]

	// Count up to 'at'.
	for i := at / r.jump*r.jump + 1; i <= at; i++ {
		if r.str[i] == char {
			result++
		}
	}

	return result
}


// ----- FIRST COLUMN INDEX ----------------------------------------------------

// An index on the first column of the BW-matrix.
type fcIndex map[byte]int

// Creates an index on the given BW-transformed sequence.
func newFcIndex(str []byte) fcIndex {
	counts := make(map[byte]int)
	for _, b := range str {
		counts[b]++
	}
	
	result := make(map[byte]int)
	for i := byte(0); i < ^byte(0); i++ {
		result[i+1] = result[i] + counts[i]
	}
	
	return result
}

// Returns the index of the given character in the first (sorted)
// column of the BW-matrix.
func (f fcIndex) indexOf(char byte, rank int) int {
	return f[char] + rank - 1  // -1 because rank is 1-based.
}
