package datum

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBuildFilteredIndex(t *testing.T) {
	content := "addr1\naddr2\naddr3\n"
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	fi, err := BuildFilteredIndex(tmpfile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, fi)
	assert.Equal(t, 3, fi.Len())

	assert.True(t, fi.Contains("addr1"))
	assert.True(t, fi.Contains("addr2"))
	assert.True(t, fi.Contains("addr3"))
	assert.False(t, fi.Contains("addr4"))
}
