package brainforce

import (
	"errors"
	"fmt"
	"sync"

	"github.com/shlima/fortune/internal/pkg/datum"
	"github.com/shlima/fortune/internal/pkg/key"
	"github.com/shlima/fortune/internal/pkg/pass"
)

type Force struct {
	index    *datum.FilteredIndex
	key      key.IGen
	pass     pass.IGen
	workers  int
	wg       *sync.WaitGroup
	mx       *sync.Mutex
	channels []PassCh
	onFound  OnFoundFn
}

func New(index *datum.FilteredIndex, key key.IGen, pass pass.IGen, workers int) *Force {
	channels := make([]PassCh, workers)
	for i := 0; i < workers; i++ {
		channels[i] = make(PassCh, 2)
	}

	return &Force{
		key:      key,
		pass:     pass,
		index:    index,
		mx:       new(sync.Mutex),
		wg:       new(sync.WaitGroup),
		workers:  workers,
		channels: channels,
	}
}

func (f *Force) SetIndex(index *datum.FilteredIndex) {
	f.index = index
}

func (f *Force) PassGen() pass.IGen {
	return f.pass
}

func (f *Force) DataLength() int {
	return f.index.Len()
}

// Get tests the index with the passed address
func (f *Force) Get(address string) bool {
	return f.index.Contains(address)
}

func (f *Force) Generate(onFound OnFoundFn) error {
	f.onFound = onFound
	f.asyncWatch()

	i := 0
LOOP:
	for {
		password, err := f.pass.Next()
		switch {
		case errors.Is(err, pass.ErrEnd):
			break LOOP
		case err != nil:
			return fmt.Errorf("failed to next: %w", err)
		}

		f.channels[i] <- password
		i = (i + 1) % f.workers
	}

	f.stop()
	return nil
}

func (f *Force) asyncWatch() {
	for i := range f.channels {
		f.wg.Add(1)
		go f.watch(f.channels[i])
	}
}

func (f *Force) stop() {
	for i := range f.channels {
		close(f.channels[i])
	}

	f.wg.Wait()
}

func (f *Force) watch(ch PassCh) {
	for password := range ch {
		chain, err := f.key.BrainSHA256([]byte(password))
		if err != nil {
			panic(fmt.Errorf("failed to key gen <%s>: %w", password, err))
		}

		if f.index.Contains(chain.Compressed) || f.index.Contains(chain.Uncompressed) {
			f.mx.Lock()
			f.onFound(chain)
			f.mx.Unlock()
		}
	}

	f.wg.Done()
}
