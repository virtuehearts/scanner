package cmd

import (
	"github.com/shlima/fortune/internal/pkg/key"
	"github.com/urfave/cli/v2"
)

func NewKeyGen(c *cli.Context) key.IGen {
	if c.Bool(FlagGPU.Name) {
		return key.NewGPU()
	}

	return key.New()
}
