package bruteforce

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shlima/fortune/internal/pkg/datum"
	"github.com/shlima/fortune/internal/pkg/key"
)

type Executor struct {
	index     *datum.FilteredIndex
	gen       key.IGen
	workers   int
	sleep     time.Duration
	ops       uint64
	prev      uint64
	wg        *sync.WaitGroup
	t0        time.Time
	mx        *sync.Mutex
	closes    []CloseCh
	batchSize int
}

func New(index *datum.FilteredIndex, gen key.IGen, workers int) *Executor {
	return &Executor{
		index:     index,
		gen:       gen,
		workers:   workers,
		wg:        new(sync.WaitGroup),
		closes:    make([]CloseCh, 0),
		t0:        time.Now(),
		mx:        new(sync.Mutex),
		batchSize: 1,
	}
}

func (e *Executor) SetBatchSize(size int) {
	if size > 0 {
		e.batchSize = size
	}
}

func (e *Executor) SetWorkers(count int) {
	e.workers = count
}

// SetNightMode sets night mode (use less CPU)
func (e *Executor) SetNightMode(on bool) {
	switch on {
	case true:
		e.sleep = time.Millisecond
	default:
		e.sleep = time.Duration(0)
	}
}

// SetIndex redefined the index
func (e *Executor) SetIndex(index *datum.FilteredIndex) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.index = index
}

// DataLength returns index items counter
func (e *Executor) DataLength() int {
	return e.index.Len()
}

func (e *Executor) WorkersCount() int {
	return e.workers
}

// Get tests the index with the passed address
func (e *Executor) Get(address string) bool {
	return e.index.Contains(address)
}

func (e *Executor) Run(fn FoundFn) {
	e.RunAsync(fn)
	e.wg.Wait()
}

func (e *Executor) RunAsync(fn FoundFn) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.wg = new(sync.WaitGroup)
	e.closes = make([]CloseCh, e.workers)
	for i := 0; i < e.workers; i++ {
		e.wg.Add(1)
		e.closes[i] = make(CloseCh, 1)
		go e.run(e.closes[i], fn)
	}
}

func (e *Executor) Stop() {
	for i := range e.closes {
		e.closes[i] <- true
	}

	e.wg.Wait()
	for i := range e.closes {
		close(e.closes[i])
	}
}

func (e *Executor) run(close CloseCh, foundFn FoundFn) {
	for {
		select {
		case <-close:
			e.wg.Done()
			return
		default:
			if e.batchSize > 1 {
				e.nextBatch(foundFn)
			} else {
				e.next(foundFn)
			}
			time.Sleep(e.sleep)
		}
	}
}

func (e *Executor) next(foundFn FoundFn) {
	pair, err := e.gen.Generate()
	if err != nil {
		panic(fmt.Errorf("failed to generate: %w", err))
	}

	if e.index.Contains(pair.Compressed) || e.index.Contains(pair.Uncompressed) {
		e.mx.Lock()
		foundFn(pair)
		e.mx.Unlock()
	}
	e.addOpts(1)
}

func (e *Executor) nextBatch(foundFn FoundFn) {
	var pairs []key.KeyChain
	var err error

	if gpuGen, ok := e.gen.(interface {
		GenerateBatch(int) ([]key.KeyChain, error)
	}); ok {
		pairs, err = gpuGen.GenerateBatch(e.batchSize)
	} else {
		pairs = make([]key.KeyChain, e.batchSize)
		for i := 0; i < e.batchSize; i++ {
			pairs[i], err = e.gen.Generate()
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		panic(fmt.Errorf("failed to generate batch: %w", err))
	}

	for _, pair := range pairs {
		if e.index.Contains(pair.Compressed) || e.index.Contains(pair.Uncompressed) {
			e.mx.Lock()
			foundFn(pair)
			e.mx.Unlock()
		}
	}
	e.addOpts(uint64(len(pairs)))
}

func (e *Executor) Heartbeat() *HeartBit {
	e.mx.Lock()
	defer e.mx.Unlock()

	t1 := time.Now()
	tried := e.GetOps()

	out := &HeartBit{
		Tried: tried,
		IOps:  uint64(math.Round(float64(tried-e.prev) / t1.Sub(e.t0).Seconds())),
	}

	e.prev = tried
	e.t0 = t1

	return out
}

func (e *Executor) GetOps() uint64 {
	return atomic.LoadUint64(&e.ops)
}

func (e *Executor) addOpts(value uint64) {
	atomic.AddUint64(&e.ops, value)
}
