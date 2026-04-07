package pass

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type FileGen struct {
	filename   string
	file       *os.File
	scanner    *bufio.Scanner
	iterations uint64
	prev       uint64
	t0         time.Time
	mx         *sync.Mutex
	lastPass   string
}

func NewFileGen(filename string) (*FileGen, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open wordlist <%s>: %w", filename, err)
	}

	return &FileGen{
		filename: filename,
		file:     file,
		scanner:  bufio.NewScanner(file),
		mx:       new(sync.Mutex),
		t0:       time.Now(),
	}, nil
}

func (g *FileGen) Next() (string, error) {
	g.mx.Lock()
	defer g.mx.Unlock()

	if g.scanner.Scan() {
		g.lastPass = strings.TrimSpace(g.scanner.Text())
		atomic.AddUint64(&g.iterations, 1)
		return g.lastPass, nil
	}

	if err := g.scanner.Err(); err != nil {
		return "", fmt.Errorf("scanner error: %w", err)
	}

	return "", ErrEnd
}

func (g *FileGen) Permutations() uint64 {
	return 0 // Unknown for file
}

func (g *FileGen) All() []string {
	return nil
}

func (g *FileGen) Heartbeat() *HeartBit {
	g.mx.Lock()
	defer g.mx.Unlock()

	t1 := time.Now()
	tried := atomic.LoadUint64(&g.iterations)

	out := &HeartBit{
		Tried:    tried,
		Password: g.lastPass,
		IOps:     uint64(tried - g.prev), // simplified for file
	}

	g.prev = tried
	g.t0 = t1

	return out
}

func (g *FileGen) Opts() GenOpts {
	return GenOpts{}
}
