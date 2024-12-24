package main

import (
	"bufio"
	"fmt"
	"os"
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
	components := CC([]int{})
	computer2components := map[string]int{}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var computer1, computer2 string
		fmt.Sscanf(scanner.Text(), "%s-%s", &computer1, &computer2)

		if _, ok := computer2components[computer1]; !ok {
			computer2components[computer1] = components.NewComponent()
		}

		if _, ok := computer2components[computer1]; !ok {
			computer2components[computer2] = components.NewComponent()
		}

		components.Union(computer2components[computer1], computer2components[computer2])
	}

	component2computers := map[int][]string{}

	for computer, component := range computer2components {
		root := components.Find(component)

		if _, exists := component2computers[root]; !exists {
			component2computers[root] = []string{}
		}

		component2computers[root] = append(component2computers[root], computer)
	}

	// we are looking for sets of size 3, where one of the computers has a t
	count := 0

	for i, computers := range component2computers {

		fmt.Printf("component %d: %d computers\n", i, len(computers))

		if len(computers) != 3 {
			continue
		}

		for _, c := range computers {
			if c[0] == 't' {
				count++
				break
			}
		}
	}

	println(count)
}
