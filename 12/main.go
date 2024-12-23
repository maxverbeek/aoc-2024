package main

import (
	"bufio"
	"os"
)

type Direction struct {
	dx, dy int
}

var (
	Top   Direction = Direction{0, -1}
	Right Direction = Direction{1, 0}
	Down  Direction = Direction{0, 1}
	Left  Direction = Direction{-1, 0}
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

func InBounds[T any](grid [][]T, x, y int) bool {
	return y >= 0 && y < len(grid) && x >= 0 && x < len(grid[y])
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

	cornerstop := [...]Direction{Top, Left, Down, Right}
	cornersleft := [...]Direction{Right, Top, Left, Down}
	// debugdirections := [...]string{"top-right", "left-top", "down-left", "right-down"}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			this := components[y][x]

			// iterate over all corner pairs, check if both corners are OOB or a different component
			for i := range cornersleft {
				cl, ct := cornersleft[i], cornerstop[i]

				checktop := InBounds(components, x+ct.dx, y+ct.dy) && components[y+ct.dy][x+ct.dx] == this
				checkleft := InBounds(components, x+cl.dx, y+cl.dy) && components[y+cl.dy][x+cl.dx] == this
				checkboth := InBounds(components, x+ct.dx+cl.dx, y+ct.dy+cl.dy) && components[y+ct.dy+cl.dy][x+ct.dx+cl.dx] == this

				if (!checktop && !checkleft) || (checktop && checkleft && !checkboth) {
					// fmt.Printf("corner for area %d at (%d, %d) in direction %s\n", this, y, x, debugdirections[i])
					corners[this]++
				}
			}
		}
	}

	part1 := 0
	part2 := 0

	for areatype, area := range area {
		part1 += area * edges[areatype]
		part2 += area * corners[areatype]
	}

	println(part1, part2)
}
