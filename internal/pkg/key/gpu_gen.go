package key

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type GPUGenerator struct {
	max     *big.Int
	testing string
	enabled bool
}

func (g *GPUGenerator) SetTesting(address string) IGen {
	g.testing = address
	return g
}

func (g *GPUGenerator) Generate() (KeyChain, error) {
	if g.testing != "" {
		return NewTestingKeyChain(g.testing), nil
	}

	res, err := rand.Int(rand.Reader, g.max)
	if err != nil {
		return KeyChain{}, fmt.Errorf("failed to rand: %w", err)
	}

	return KeyChainFromPriv(res.Bytes())
}

func (g *GPUGenerator) BrainSHA256(password []byte) (KeyChain, error) {
	return KeyChain{}, fmt.Errorf("GPU generator does not support BrainSHA256")
}

func (g *GPUGenerator) generateBatchCPU(count int) ([]KeyChain, error) {
	out := make([]KeyChain, count)
	for i := 0; i < count; i++ {
		key, err := g.Generate()
		if err != nil {
			return nil, err
		}
		out[i] = key
	}
	return out, nil
}
