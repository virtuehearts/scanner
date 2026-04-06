package datum

import (
	"math/rand"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
)

type Index map[string]bool

type FilteredIndex struct {
	index  Index
	filter *bloom.BloomFilter
}

func NewFilteredIndex(index Index, filter *bloom.BloomFilter) *FilteredIndex {
	return &FilteredIndex{
		index:  index,
		filter: filter,
	}
}

func (i *FilteredIndex) Contains(address string) bool {
	if i.filter != nil && !i.filter.TestString(address) {
		return false
	}

	return i.index[address]
}

func (i *FilteredIndex) Len() int {
	return len(i.index)
}

func (i *FilteredIndex) Random() (address string, ix int) {
	return i.index.Random()
}

// Random returns random address from the index
func (i Index) Random() (address string, ix int) {
	if len(i) == 0 {
		return
	}

	number := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(i))

LOOP:
	for key, _ := range i {
		switch number {
		case ix:
			address = key
			break LOOP
		default:
			ix++
		}
	}

	return
}

func (i Index) SetTesting(address string) Index {
	if address == "" {
		return i
	}

	return Index{address: true}
}

func (i *FilteredIndex) SetTesting(address string) *FilteredIndex {
	if address == "" {
		return i
	}

	i.index = i.index.SetTesting(address)
	return i
}
