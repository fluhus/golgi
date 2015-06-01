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

