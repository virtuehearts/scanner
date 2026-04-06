package key

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGPUGenerator_Generate(t *testing.T) {
	g := NewGPU()
	key, err := g.Generate()
	assert.NoError(t, err)
	assert.NotEmpty(t, key.Private)
	assert.NotEmpty(t, key.Compressed)
	assert.NotEmpty(t, key.Uncompressed)
}

func TestGPUGenerator_GenerateBatch(t *testing.T) {
	g := NewGPU()
	keys, err := g.GenerateBatch(10)
	assert.NoError(t, err)
	assert.Len(t, keys, 10)
	for _, key := range keys {
		assert.NotEmpty(t, key.Private)
		assert.NotEmpty(t, key.Compressed)
		assert.NotEmpty(t, key.Uncompressed)
	}
}
