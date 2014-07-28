package main

import (
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/zhirsch/concurrentrand"
)

var (
	startSize = flag.Int("start_size", 1, "the amount of numbers to start with")
)

func worker(n []int, r *rand.Rand, trych <-chan struct{}, donech chan<- struct{}) {
	m := numbers(make([]int, len(n)))
	copy(m, n)
	for _ = range trych {
		m.shuffle(r)
		if m.sorted() {
			donech <- struct{}{}
			return
		}
	}
}

func report(start time.Time, size int, count int64) {
	dur := time.Since(start)
	hertz := float64(count) / dur.Seconds()
	per := time.Duration(dur.Nanoseconds()/count) * time.Nanosecond
	fmt.Printf("Size:      %d\n", size)
	fmt.Printf("Tries:     %d\n", count)
	fmt.Printf("Elapsed:   %v\n", dur)
	fmt.Printf("Tries/sec: %.f\n", hertz)
	fmt.Printf("Secs/try:  %v\n", per)
	fmt.Println()
}

func sort(n []int, workers int) {
	trych := make(chan struct{}, workers*500)
	donech := make(chan struct{})

	// Start all the workers.
	rsrc := concurrentrand.NewSource(1)
	for i := 0; i < workers; i++ {
		go worker(n, rand.New(rsrc), trych, donech)
	}

	var count int64
	start := time.Now()

	reporter := func() { report(start, len(n), count) }
	defer reporter()

	progressch := time.Tick(15 * time.Second)

	// Loop until one of the workers sorts the list.
	for {
		select {
		case <-donech:
			// Sorted, all done.
			close(trych)
			return
		case trych <- struct{}{}:
			// More work is needed.
			count += 1
		case <-progressch:
			// Periodically report progress.
			reporter()
		}
	}
}

func main() {
	cpu := runtime.NumCPU()
	fmt.Printf("GOMAXPROCS=%d\n", cpu)
	runtime.GOMAXPROCS(cpu)

	size := *startSize
	for {
		var n []int
		for i := 0; i < size; i++ {
			n = append(n, i)
		}
		sort(n, cpu)
		size++
	}
}
