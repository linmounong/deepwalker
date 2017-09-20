package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

const (
	logEvery = 50000
)

type Graph map[uint32][]uint32

func contains(a []uint32, i uint32) bool {
	for _, v := range a {
		if v == i {
			return true
		}
	}
	return false
}

func (g Graph) RandWalk(id uint32) uint32 {
	a := g[id]
	return a[rand.Intn(len(a))]
}

func (g Graph) AddNode(id uint32) {
	if _, ok := g[id]; !ok {
		g[id] = []uint32{}
		if len(g)%logEvery == 0 {
			log.Println("len(g)", len(g), "still loading...")
		}
	}
}

func (g Graph) Equals(g2 Graph) bool {
	if len(g) != len(g2) {
		return false
	}
	for k, v := range g {
		v2, ok := g2[k]
		if !ok || len(v) != len(v2) {
			return false
		}
		for _, vv := range v {
			if !contains(v2, vv) {
				return false
			}
		}
	}
	return true
}

func (g Graph) AddEdge(from uint32, tolist ...uint32) {
	g.AddNode(from)
	for _, to := range tolist {
		g.AddNode(to)
		if !contains(g[from], to) {
			g[from] = append(g[from], to)
		}
		if !contains(g[to], from) {
			g[to] = append(g[to], from)
		}
	}
}

func (g Graph) DebugString() string {
	var buf bytes.Buffer
	for k, v := range g {
		buf.WriteString(fmt.Sprintf("%v -> %v\n", k, v))
	}
	return buf.String()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func NewGraphFromAdjlist(reader io.Reader) Graph {
	g := Graph{}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer([]byte{}, 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var from uint32
		for i, col := range strings.Split(line, " ") {
			to, err := strconv.ParseUint(col, 10, 32)
			check(err)
			if i == 0 {
				from = uint32(to)
			} else {
				g.AddEdge(from, uint32(to))
			}
		}
	}
	check(scanner.Err())
	return g
}
