package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseURI(t *testing.T) {
	tests := []struct {
		valid           bool
		input, expected string
	}{
		{true, ":2444", "http://localhost:2444"},
	}

	for _, test := range tests {
		uri, err := baseURI(test.input)
		if !test.valid {
			assert.Error(t, err)
		}
		assert.Equal(t, test.expected, uri)
	}
}
