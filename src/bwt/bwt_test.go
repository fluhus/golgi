package bwt

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTransform_word(t *testing.T) {
	assert := assert.New(t)
	out := string(Transform([]byte("abracadabra")))
	assert.Equal("ard$rcaaaabb", out)
}

func TestTransform_empty(t *testing.T) {
	assert := assert.New(t)
	out := string(Transform(nil))
	assert.Equal("$", out)
}

func TestRankIndex_simple(t *testing.T) {
	assert := assert.New(t)
	in := "aaaa"
	idx := newRankIndex([]byte(in), 1)
	for i := range in {
		assert.Equal(i+1, idx.rankOf('a', i))
		assert.Equal(0, idx.rankOf('b', i))
	}
}

func TestRankIndex_complex(t *testing.T) {
	assert := assert.New(t)
	in := "aabbaaabcccb"
	for jump := 1; jump <= 20; jump++ {
		idx := newRankIndex([]byte(in), jump)
		assert.Equal(1, idx.rankOf('a', 0))
		assert.Equal(2, idx.rankOf('a', 1))
		assert.Equal(4, idx.rankOf('a', 5))
		assert.Equal(5, idx.rankOf('a', 6))
		assert.Equal(5, idx.rankOf('a', 7))
		assert.Equal(5, idx.rankOf('a', 11))

		assert.Equal(0, idx.rankOf('b', 1))
		assert.Equal(1, idx.rankOf('b', 2))
		assert.Equal(2, idx.rankOf('b', 5))
		assert.Equal(3, idx.rankOf('b', 10))
		assert.Equal(4, idx.rankOf('b', 11))

		assert.Equal(0, idx.rankOf('c', 7))
		assert.Equal(1, idx.rankOf('c', 8))
		assert.Equal(3, idx.rankOf('c', 11))
	}
}







