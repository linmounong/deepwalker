package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

var (
	input       = flag.String("input", "-", "input graph file in adjlist format")
	output      = flag.String("output", "-", "output representation file")
	workers     = flag.Int("workers", 1, "number of parallel go routines")
	numberWalks = flag.Int("number-walks", 10, "number of random walks to start at each node")
	saveVocab   = flag.String("save-vocab", "", "file to save vocab, format conforming to the c implementation")
)

const (
	flushEvery = 50000
)

func dumpVocab(vocab map[uint32]uint64, out *bufio.Writer) {
	out.WriteString("</s> 0\n") // magic
	type kv struct {
		id  uint32
		cnt uint64
	}
	sorted := make([]kv, 0, len(vocab))
	for id, cnt := range vocab {
		sorted = append(sorted, kv{id, cnt})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].cnt > sorted[j].cnt
	})
	for _, item := range sorted {
		out.WriteString(fmt.Sprintf("%d %d\n", item.id, item.cnt))
	}
	out.Flush()
}

func main() {
	flag.Parse()
	var reader io.Reader
	if *input == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(*input)
		check(err)
		defer f.Close()
		reader = bufio.NewReader(f)
	}
	var out *bufio.Writer
	if *output == "-" {
		out = bufio.NewWriter(os.Stdout)
	} else {
		f, err := os.Create(*output)
		check(err)
		defer f.Close()
		out = bufio.NewWriter(f)
	}
	defer out.Flush()
	g := NewGraphFromAdjlist(reader)
	log.Println("loaded g, len(g)", len(g))
	walksPool := []<-chan []uint32{}
	ids := make(chan uint32)
	for i := 0; i < *workers; i++ {
		walks := make(chan []uint32)
		walksPool = append(walksPool, walks)
		go Walk(g, ids, walks)
	}
	go func() {
		for i := 0; i < *numberWalks; i++ {
			for k := range g {
				ids <- k
			}
		}
		close(ids)
	}()
	var cnt int
	var vocab map[uint32]uint64
	if *saveVocab != "" {
		vocab = map[uint32]uint64{}
	}
	for walk := range MergeWalks(walksPool...) {
		size := len(walk)
		for i, id := range walk {
			_, err := out.WriteString(strconv.FormatUint(uint64(id), 10))
			check(err)
			if i < size {
				_, err := out.WriteString(" ")
				check(err)
			}
			if vocab != nil {
				vocab[id] += 1
			}
		}
		_, err := out.WriteString("\n")
		check(err)
		cnt += 1
		if cnt%flushEvery == 0 {
			log.Println("wrote", cnt, "lines")
			out.Flush()
		}
	}
	if vocab != nil {
		f, err := os.Create(*saveVocab)
		check(err)
		defer f.Close()
		out = bufio.NewWriter(f)
		dumpVocab(vocab, out)
	}
}
