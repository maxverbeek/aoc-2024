package main

import (
	"bufio"
	"os"
)

func Find(parents []int, id int) int {
	root := parents[id]

	for root != parents[root] {
		root = parents[root]
	}

	parents[id] = root

	return root
}

func Unify(parents []int, a, b int) {
	a, b = Find(parents, a), Find(parents, b)

	if a < b {
		parents[b] = a
	} else if a > b {
		parents[a] = b
	}
}

func NewComponent(parents *[]int) int {
	newcomponent := len(*parents)
	*parents = append(*parents, newcomponent)

	return newcomponent
}

func main() {
	world := [][]rune{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		world = append(world, []rune(scanner.Text()))
	}

	height := len(world)
	width := len(world[0])

	parents := []int{}
	components := make([][]int, height)

	for y := range components {
		components[y] = make([]int, width)
	}

	// union find algorithm to find connected components
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			this := world[y][x]
			components[y][x] = NewComponent(&parents)

			if y > 0 && this == world[y-1][x] {
				Unify(parents, components[y][x], components[y-1][x])
			}

			if x > 0 && this == world[y][x-1] {
				Unify(parents, components[y][x], components[y][x-1])
			}
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			components[y][x] = Find(parents, components[y][x])
		}
	}

	edges := map[int]int{}
	area := map[int]int{}
	corners := map[int]int{}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// count the surface as well
			area[components[y][x]]++

			if x == 0 || x == width-1 {
				// add an edge around the boundary of the map
				edges[components[y][x]]++
			}

			if x < width-1 && world[y][x] != world[y][x+1] {
				edges[components[y][x]]++
				edges[components[y][x+1]]++
			}
		}
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if y == 0 || y == height-1 {
				// add an edge around the boundary of the map
				edges[components[y][x]]++
			}

			if y < height-1 && world[y][x] != world[y+1][x] {
				edges[components[y][x]]++
				edges[components[y+1][x]]++
			}
		}
	}

	// count corners (this is wrong)
	for y := 0; y < height-1; y++ {
		for x := 0; x < width; x++ {
			if x-1 >= 0 && components[y][x] != components[y+1][x-1] {
				corners[components[y][x]]++
				corners[components[y+1][x-1]]++
			}

			if x+1 < width && components[y][x] != components[y+1][x+1] {
				corners[components[y][x]]++
				corners[components[y+1][x+1]]++
			}
		}
	}

	part1 := 0
	part2 := 0

	for areatype, area := range area {
		perimeter := edges[areatype]
		part1 += area * perimeter
	}

	for areatype, area := range area {
		numedges := corners[areatype]
		part2 += area * numedges
	}

	println(part1, part2)
}
