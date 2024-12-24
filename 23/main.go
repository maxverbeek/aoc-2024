package main

import (
	"bufio"
	"maps"
	"os"
	"slices"
	"strings"
)

// connected components
type CC []int

func (cc *CC) NewComponent() int {
	newcomp := len(*cc)
	*cc = append(*cc, newcomp)

	return newcomp
}

func (cc *CC) Find(component int) int {
	// Find the parent by iterating over the parent array. If parents[i] == i then we are at the root.

	root := component

	for []int(*cc)[root] != root {
		root = []int(*cc)[root]
	}

	[]int(*cc)[component] = root

	return root
}

func (cc *CC) Union(a, b int) {
	a, b = cc.Find(a), cc.Find(b)

	if a < b {
		[]int(*cc)[b] = a
	} else {
		[]int(*cc)[a] = b
	}
}

func main() {
	graph := map[string]map[string]bool{}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		connection := strings.Split(scanner.Text(), "-")

		computer1, computer2 := connection[0], connection[1]

		if _, exists := graph[computer1]; !exists {
			graph[computer1] = map[string]bool{}
		}

		if _, exists := graph[computer2]; !exists {
			graph[computer2] = map[string]bool{}
		}

		graph[computer1][computer2] = true
		graph[computer2][computer1] = true
	}

	computers := slices.Sorted(maps.Keys(graph))

	triples := [][3]string{}

	for i, c1 := range computers {
		for j, c2 := range computers[i+1:] {
			for _, c3 := range computers[i+j+1:] {
				if !graph[c1][c2] || !graph[c2][c3] || !graph[c3][c1] {
					continue
				}

				if c1[0] == 't' || c2[0] == 't' || c3[0] == 't' {
					triples = append(triples, [...]string{c1, c2, c3})
				}
			}
		}
	}

	println(len(triples))

	mostconnections := 0
	mostconnected := []string{}

	for comp, edges := range graph {
		if len(edges) > mostconnections {
			mostconnections = len(edges)
			mostconnected = []string{comp}
		} else if len(edges) == mostconnections {
			mostconnected = append(mostconnected, comp)
		}
	}

	slices.Sort(mostconnected)

	println(strings.Join(mostconnected, ","))
}
