package cmd

import (
	"fmt"

	"github.com/shlima/fortune/internal/pkg/datum"
	"github.com/urfave/cli/v2"
)

func NewIndex(c *cli.Context) *datum.FilteredIndex {
	var index *datum.FilteredIndex
	var err error

	if c.Bool(FlagBloom.Name) {
		index, err = datum.BuildFilteredIndex(c.StringSlice(FlagFiles.Name)...)
	} else {
		idx, readErr := datum.ReadFiles(c.StringSlice(FlagFiles.Name)...)
		if readErr == nil {
			index = datum.NewFilteredIndex(idx, nil)
		}
		err = readErr
	}

	if err != nil {
		panic(fmt.Errorf("failed to read datum: %w", err))
	}

	return index
}
