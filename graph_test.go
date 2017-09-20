package main

import (
	"strings"
	"testing"
)

func TestGraph(t *testing.T) {
	adjlist := `
1 2 3
2 3
3 4
4 5`
	reader := strings.NewReader(adjlist)
	g2 := NewGraphFromAdjlist(reader)

	g1 := Graph{}
	g1.AddEdge(1, 2, 3)
	g1.AddEdge(2, 3)
	g1.AddEdge(3, 4)
	if g1.Equals(g2) {
		t.Errorf("g1 == g2\n\n%s\n\n%s", g1.DebugString(), g2.DebugString())
	}
	g1.AddEdge(4, 5)
	if !g1.Equals(g2) {
		t.Errorf("g1 != g2\n\n%s\n\n%s", g1.DebugString(), g2.DebugString())
	}
	// test duplicate
	g1.AddEdge(4, 5)
	if !g1.Equals(g2) {
		t.Errorf("g1 != g2\n\n%s\n\n%s", g1.DebugString(), g2.DebugString())
	}
}
