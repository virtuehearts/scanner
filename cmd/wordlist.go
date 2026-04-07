package cmd

import (
	"fmt"

	"github.com/shlima/fortune/internal/pkg/brainforce"
	"github.com/shlima/fortune/internal/pkg/pass"
	"github.com/urfave/cli/v2"
)

func Wordlist(c *cli.Context) error {
	filename := c.String(FlagPassFile.Name)
	if filename == "" {
		return fmt.Errorf("wordlist file is required")
	}

	p, err := pass.NewFileGen(filename)
	if err != nil {
		return fmt.Errorf("failed to create file gen: %w", err)
	}

	force := brainforce.New(
		NewIndex(c).SetTesting(c.Args().First()),
		NewKeyGen(c),
		p,
		c.Int(FlagWorkers.Name),
	)

	go brainForceHeartbit(c, force)
	go brainForceTelegram(c, force)
	return force.Generate(onFound(c))
}
