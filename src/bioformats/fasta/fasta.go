// Deals with fasta parsing and representation.
package fasta

import (
	//"os"
	"fmt"
	//"bufio"
	//"strings"
)

// Converts number to nucleotide.
var num2nuc = []byte{'A', 'C', 'G', 'T'}

// Converts nucleotide to number.
var nuc2num = map[byte]uint {'A':0, 'C':1, 'G':2, 'T':3, 'N':4, 'a':0, 'c':1, 'g':2, 't':3, 'n':4}


// *** FASTA ENTRY ************************************************************

// A single fasta sequence.
type FastaEntry struct {
	name     string            // sequence name (row that starts with '>')
	sequence []byte            // sequence in 2-bit format
	length   uint              // number of nucleotides
	isN      map[uint]struct{} // coordinates of 'N' nucleotides
}

// Returns an empty fasta entry.
func newFastaEntry() *FastaEntry {
	return &FastaEntry{"", nil, 0, make(map[uint]struct{})}
}

// Returns the number of nucleotides in this fasta entry.
func (f *FastaEntry) Length() int {
	return int(f.length)
}

// Returns the name of the fasta entry.
func (f *FastaEntry) Name() string {
	return f.name
}

// Returns the nucleotide at the given position.
func (f *FastaEntry) At(position int) byte {
	uposition := uint(position)

	// Check if N
	if _,n := f.isN[uposition]; n {
		return 'N'
	}

	// Extract nucleotide
	num := (f.sequence[uposition / 4] >> (uposition % 4 * 2) & 3)
	
	return num2nuc[num]
}

// Appends a nucleotide to the fasta entry.
func (f *FastaEntry) append(nuc byte) {
	num, ok := nuc2num[nuc]

	// If unknown nucleotide
	if !ok {
		panic("Bad nucleotide: " + string([]byte{nuc}))
	}

	// If 'N'
	if num == 4 {
		num = 0
		f.isN[f.length] = struct{}{}
	}

	// Append an extra byte
	if f.length % 4 == 0 {
		f.sequence = append(f.sequence, 0)
	}

	// Set bits
	f.sequence[f.length / 4] |= byte( num << (f.length % 4 * 2) )

	f.length++
}

// String representation of an entry. Format: name[length]
func (f *FastaEntry) String() string {
	return fmt.Sprintf("%s[%d]", f.name, f.Length())
}


// *** FASTA ******************************************************************

/*
// An entire fasta file.
// Contains an array of fasta sequences.
type Fasta []FastaEntry

// String representation of a fasta file.
func (f Fasta) String() string {
	result := ""
	for i,v := range f {
		result += fmt.Sprintf("(%d)\t%s\n", i, v.String())
	}
	return result
}

// Reads fasta from a file. Returns nil on error.
func FastaFromFile(path string) (Fasta, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {return nil, err}
	defer file.Close()

	// Create buffer (1M buffer size)
	reader := bufio.NewReader(file)

	// Start reading
	result := Fasta(nil)
	first := true
	for l, e := reader.ReadBytes('\n'); e == nil; l,
		e = reader.ReadBytes('\n') {
		// Trim line delimiters
		l = []byte(strings.Trim(string(l), "\r\n"))
		
		// Ignore zero length line
		if len(l) == 0 {continue}

		// First line - no need for a title
		if first && l[0] != '>' {
			// Add first entry
			result = append(result, FastaEntry{[]byte("(no title)"), nil})
		}

		// Not first anymore
		first = false

		// If starting with '>' start a new entry
		if l[0] == '>' {
			result = append(result, FastaEntry{l[1:], nil})

		// If not, add to current sequence
		} else {
			result[len(result) - 1].Sequence =
					append(result[len(result) - 1].Sequence, l...)
		}
	}

	return result, nil
}

// Returns the index of the first sequence, whose title equals the given one.
// Returns -1 if no such sequence was found.
func (f Fasta) indexOf(title string) int {
	for i := range f {
		if string(f[i].Title) == title {
			return i
		}
	}

	// Nothing was found =[
	return -1
}
// BUG( ) Find a better name for 'indexOf'.

// Counts how many (overlapping) subsequences of the given
// length exist in the fasta.
func (f Fasta) NumberOfSubsequences(length int) int {
	// Check input
	if length < 1 {
		panic(fmt.Sprint("bad length: ", length))
	}

	result := 0

	// Count subsequences for each chromosome
	for _,chrom := range f {
		// Number of subsequences in this chromosome
		n := len(chrom.Sequence) - length + 1
		if n < 0 { n = 0 }   // if sequence is shorter than length
		
		result += n
	}
	
	return result
}

// Returns the subsequence with the given length, at the given serial index.
// Also returns the chromosome and position of the subsequence.
func (f Fasta) Subsequence(length int, index int) (subseq []byte,
		chrom int, pos int) {
	relativeIndex := index

	for i := range f {
		// Number of subsequences in this chromosome
		n := len(f[i].Sequence) - length + 1
		if n < 0 { n = 0 }   // if sequence is shorter than length
		
		// If in this chromosome, return it
		if relativeIndex < n {
			return f[i].Sequence[relativeIndex : relativeIndex+length],
					i, relativeIndex
			
		// Else, advance index to next chromosome
		} else {
			relativeIndex -= n
		}
	}
	
	// Index is to high :(
	panic(fmt.Sprint("bad subsequence index: ", index,
			" (only ", f.NumberOfSubsequences(length), " exist)"))
}

//*/



