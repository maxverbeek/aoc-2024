package main

import (
	"bufio"
	"maps"
	"os"
	"slices"
	"strings"
)

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Delete(item T) {
	delete(s, item)
}

func (s Set[T]) Contains(item T) bool {
	_, exists := s[item]

	return exists
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	result := s.Clone()
	maps.Copy(result, other)

	return result
}

func (s Set[T]) Intersect(other Set[T]) Set[T] {
	result := Set[T]{}
	for item := range s {
		if other.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

func (s Set[T]) Remove(other Set[T]) Set[T] {
	result := s.Clone()

	for key := range other {
		delete(result, key)
	}

	return result
}

func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

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

	vertices := Set[string]{}

	for v := range graph {
		vertices.Add(v)
	}

	fcsubgraphs := BronKerbosch(graph, Set[string]{}, vertices, Set[string]{})

	largest := fcsubgraphs[0]

	for _, sg := range fcsubgraphs[1:] {
		if len(sg) > len(largest) {
			largest = sg
		}
	}

	println(strings.Join(slices.Sorted(maps.Keys(largest)), ","))
}

func BronKerbosch(graph map[string]map[string]bool, R, P, X Set[string]) []Set[string] {
	if len(P) == 0 && len(X) == 0 {
		return []Set[string]{R.Clone()}
	}

	R, P, X = R.Clone(), P.Clone(), X.Clone()

	result := []Set[string]{}

	for vertex := range P {
		neighbours := Set[string]{}
		for neighbour := range graph[vertex] {
			neighbours.Add(neighbour)
		}

		R.Add(vertex)
		result = append(result, BronKerbosch(graph, R, P.Intersect(neighbours), X.Intersect(neighbours))...)
		delete(R, vertex)

		delete(P, vertex)
		X.Add(vertex)
	}

	return result
}

func RemoveFromSlice[T any](slice []T, idx int) []T {
	last := len(slice) - 1
	slice[idx] = slice[last]
	return slice[:last]
}
