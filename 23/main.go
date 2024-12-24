package main

import (
	"bufio"
	"maps"
	"os"
	"reflect"
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

	// start out with fully connected singleton graphs
	fcsubgraphs := []map[string]struct{}{}
	for node := range graph {
		fcsubgraphs = append(fcsubgraphs, map[string]struct{}{node: struct{}{}})
	}

	// iterate over all edges of the graph, and check for each edge if the
	// connected node not yet part of the graph is connected to every node in
	// the graph

	degree := 1

	for {
		// increase each subgraph by 1 degree
		nextdegreesubgraphs := []map[string]struct{}{}
		for _, sg := range fcsubgraphs {
			// keep track of neighbours in a set so we don't consider duplicates
			neighbours := map[string]struct{}{}

			for node := range sg {
				for neighbour := range graph[node] {
					if len(graph[neighbour]) <= degree {
						// this neighbour does not have enough edges to be a reasonable candidate
						continue
					}
					neighbours[neighbour] = struct{}{}
				}
			}

			// for each neighbour, check if they are connected to every node in
			// the graph. if so we have a new subgraph of SG + [neighbour]
			for neighbour := range neighbours {
				connected := true
				for node := range sg {
					if !graph[node][neighbour] {
						connected = false
						break
					}
				}

				if connected {
					sg2 := maps.Clone(sg)
					sg2[neighbour] = struct{}{}
					nextdegreesubgraphs = append(nextdegreesubgraphs, sg2)
				}
			}
		}

		// prune next degree subgraphs, because there should be duplicates when
		// 2 different graphs merge
		for i := len(nextdegreesubgraphs) - 1; i >= 0; i-- {
			for j := range nextdegreesubgraphs[:len(nextdegreesubgraphs)-i-1] {
				if reflect.DeepEqual(nextdegreesubgraphs[i], nextdegreesubgraphs[j]) {
					// swap with last element of the list
					nextdegreesubgraphs = RemoveFromSlice(nextdegreesubgraphs, i)
					break
				}
			}
		}

		if len(nextdegreesubgraphs) == 0 {
			break
		}

		fcsubgraphs = nextdegreesubgraphs
		degree++

		println(degree)
	}

	println(strings.Join(slices.Sorted(maps.Keys(fcsubgraphs[0])), ","))
}

func RemoveFromSlice[T any](slice []T, idx int) []T {
	last := len(slice) - 1
	slice[idx] = slice[last]
	return slice[:last]
}
