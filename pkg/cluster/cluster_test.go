package cluster

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchFilter(t *testing.T) {
	const refused = "ssh: connect to host 172.17.0.2 port 22: Connection refused"

	filter := matchFilter{
		writer: ioutil.Discard,
		regexp: connectRefused,
	}

	filter.Write([]byte("foo\n"))
	assert.Equal(t, false, filter.matched)

	filter.Write([]byte(refused))
	assert.Equal(t, true, filter.matched)
}
