//go:build !gpu
// +build !gpu

package key

import (
	"fmt"
	"math/big"
	"os"
)

func NewGPU() *GPUGenerator {
	g := &GPUGenerator{
		max: big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil),
	}

	fmt.Fprintf(os.Stderr, "GPU Generator initialized in CPU-fallback mode. Rebuild with '-tags gpu' to enable OpenCL.\n")
	g.enabled = false

	return g
}

func (g *GPUGenerator) GenerateBatch(count int) ([]KeyChain, error) {
	return g.generateBatchCPU(count)
}
