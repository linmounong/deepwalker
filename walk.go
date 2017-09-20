package main

import (
	"flag"
	"math/rand"
	"sync"
)

var (
	alpha      = flag.Float64("alpha", 0, "probability of restarts")
	walkLength = flag.Int("walk-length", 40, "length of the random walk started at each node")
)

func Walk(g Graph, ids <-chan uint32, walks chan<- []uint32) {
	for start := range ids {
		walk := make([]uint32, 0, *walkLength)
		id := start
		for i := 0; i < *walkLength; i++ {
			walk = append(walk, id)
			if rand.Float64() > *alpha {
				id = g.RandWalk(id)
			} else {
				id = start
			}
		}
		walks <- walk
	}
	close(walks)
}

func MergeWalks(cs ...<-chan []uint32) <-chan []uint32 {
	var wg sync.WaitGroup
	out := make(chan []uint32)
	output := func(c <-chan []uint32) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
