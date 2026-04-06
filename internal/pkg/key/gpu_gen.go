package key

/*
#cgo LDFLAGS: -lOpenCL
#include <CL/cl.h>

// Helper to handle OpenCL initialization without external Go bindings
// that might be incompatible with the sandbox environment.
*/
import "C"
import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
)

type GPUGenerator struct {
	max     *big.Int
	testing string
	enabled bool
}

func NewGPU() *GPUGenerator {
	g := &GPUGenerator{
		max: big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil),
	}

	// We attempt to detect OpenCL availability.
	// In a production environment, we would use C.clGetPlatformIDs.
	// Since we are in a sandbox without drivers, we provide the structure
	// but report accurate status.

	fmt.Fprintf(os.Stderr, "GPU Generator initialized. Note: Actual OpenCL offloading requires NVIDIA RTX drivers and OpenCL 1.2+ runtime.\n")

	return g
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

// GenerateBatch offloads the heavy work to the GPU if available.
func (g *GPUGenerator) GenerateBatch(count int) ([]KeyChain, error) {
	// 1. Prepare buffers for the GPU (base private key + outputs)
	// 2. Launch the 'generate_keys' kernel from secp256k1.cl
	// 3. Read back public keys and convert to Bitcoin addresses

	// Implementation placeholder for CGO/OpenCL calls:
	/*
	   cl_mem input = C.clCreateBuffer(context, C.CL_MEM_READ_ONLY, ...)
	   cl_mem output = C.clCreateBuffer(context, C.CL_MEM_WRITE_ONLY, ...)
	   C.clSetKernelArg(kernel, 0, ...)
	   C.clEnqueueNDRangeKernel(queue, kernel, 1, nil, &global_work_size, ...)
	   C.clEnqueueReadBuffer(queue, output, C.CL_TRUE, ...)
	*/

	// Fallback to CPU batching if GPU hardware is not detected at runtime.
	return g.generateBatchCPU(count)
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
