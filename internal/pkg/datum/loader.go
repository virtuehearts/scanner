package datum

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bits-and-blooms/bloom/v3"
)

func ReadFile(filename string) (Index, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file <%s>: %w", filename, err)
	}
	defer file.Close()

	index := make(Index)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			index[line] = true
		}
	}

	return index, nil
}

func ReadFiles(filenames ...string) (Index, error) {
	index := make(Index)
	for _, name := range filenames {
		i, err := ReadFile(name)
		if err != nil {
			return nil, err
		}

		for k := range i {
			index[k] = true
		}
	}

	return index, nil
}

func BuildFilteredIndex(filenames ...string) (*FilteredIndex, error) {
	index, err := ReadFiles(filenames...)
	if err != nil {
		return nil, err
	}

	n := uint(len(index))
	if n == 0 {
		return NewFilteredIndex(index, nil), nil
	}

	filter := bloom.NewWithEstimates(n, 0.001)
	for k := range index {
		filter.AddString(k)
	}

	return NewFilteredIndex(index, filter), nil
}
