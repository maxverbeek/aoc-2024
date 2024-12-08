package main

import (
	"bufio"
	"fmt"
	"iter"
	"os"
)

type Location struct {
	x, y int
}

func (l Location) String() string {
	return fmt.Sprintf("(%d %d)", l.x, l.y)
}

func (loc1 Location) ResonateWith(loc2 Location) (Location, Location, Location, Location) {
	x1, x2 := Resonate(loc1.x, loc2.x)
	y1, y2 := Resonate(loc1.y, loc2.y)

	// also return the original locations so it's easier to iterate over the range of all resonant locations
	// by doing something like a, b, _, _ = ResonateWith(...)
	return Location{x1, y1}, loc1, loc2, Location{x2, y2}
}

func Resonate(a, b int) (int, int) {
	diff := b - a
	return a - diff, b + diff
}

// Golang 1.23 iterator abuse where I don't return an index in the first value???
func Permutations2[T any](things []T) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		for i, item1 := range things {
			for _, item2 := range things[i+1:] {
				if !yield(item1, item2) {
					return
				}
			}
		}
	}
}

func InBounds(world [][]rune, loc Location) bool {
	return loc.y >= 0 && loc.y < len(world) && loc.x >= 0 && loc.x < len(world[loc.y])
}

func Print(world [][]rune, antinodes map[Location]bool) {
	for y, worldy := range world {
		for x, rune := range worldy {
			if rune != '.' {
				fmt.Printf("%c", rune)
			} else if antinodes[Location{x, y}] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}

		fmt.Printf("\n")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	world := [][]rune{}

	for scanner.Scan() {
		world = append(world, []rune(scanner.Text()))
	}

	antennas := map[rune][]Location{}

	for y, worldy := range world {
		for x, r := range worldy {
			if r == '.' {
				continue
			}

			if _, ok := antennas[r]; !ok {
				antennas[r] = []Location{{x, y}}
			} else {
				antennas[r] = append(antennas[r], Location{x, y})
			}
		}
	}

	part1 := map[Location]bool{}
	part2 := map[Location]bool{}

	for _, locations := range antennas {
		for original1, original2 := range Permutations2(locations) {
			// part 1
			a, _, _, b := original1.ResonateWith(original2)

			if InBounds(world, a) {
				part1[a] = true
			}

			if InBounds(world, b) {
				part1[b] = true
			}

			// part 2
			loc1, loc2 := original1, original2
			for InBounds(world, loc1) {
				part2[loc1] = true
				loc1, loc2, _, _ = loc1.ResonateWith(loc2)
			}

			loc1, loc2 = original1, original2
			for InBounds(world, loc1) {
				part2[loc1] = true
				_, _, loc1, loc2 = loc1.ResonateWith(loc2)
			}
		}
	}

	// Print(world, part2)
	println(len(part1), len(part2))
}
