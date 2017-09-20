package main

import (
	"testing"
)

func TestWalk(t *testing.T) {
	g := Graph{}
	g.AddEdge(1, 2, 3)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)

	ids := make(chan uint32)
	walksPool := []<-chan []uint32{}
	for i := 0; i < 4; i++ {
		walks := make(chan []uint32)
		walksPool = append(walksPool, walks)
		go Walk(g, ids, walks)
	}
	go func() {
		for k := range g {
			ids <- k
		}
		for k := range g {
			ids <- k
		}
		for k := range g {
			ids <- k
		}
		close(ids)
	}()
	cnt := 0
	for walk := range MergeWalks(walksPool...) {
		cnt += 1
		if len(walk) != *walkLength {
			t.Errorf("len(walk) should be %d, got %d, %v", *walkLength, len(walk), walk)
		}
	}
	if cnt != len(g)*3 {
		t.Error("too few walks")
	}
}
