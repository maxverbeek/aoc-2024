package main

import (
	"bufio"
	"os"
)

type Direction int

const (
	DirUp Direction = iota
	DirRight
	DirDown
	DirLeft
)

func (d Direction) TurnRight() Direction {
	return (d + 1) % 4
}

func (d Direction) MoveIndex(x, y int) (int, int) {
	if d == DirUp {
		y--
	}

	if d == DirRight {
		x++
	}

	if d == DirDown {
		y++
	}

	if d == DirLeft {
		x--
	}

	return x, y
}

func Print(world [][]rune) {
	for _, line := range world {
		println(string(line))
	}
}

type SeenSet [][]int

func NewSeenSetFor(world [][]rune) SeenSet {
	seen := [][]int{}

	for _, line := range world {
		seen = append(seen, make([]int, len(line)))
	}

	return SeenSet(seen)
}

func (s SeenSet) Visit(x, y int, d Direction) {
	s[y][x] |= 1 << d
}

func (s SeenSet) SeenDirection(x, y int, d Direction) bool {
	return s[y][x]&(1<<d) > 0
}

func (s SeenSet) SeenAny(x, y int) bool {
	return s[y][x] > 0
}

func CountSteps(world [][]rune, locx, locy int) int {
	direction := DirUp

	seen := NewSeenSetFor(world)
	seen.Visit(locx, locy, direction)

	for {
		newx, newy := direction.MoveIndex(locx, locy)

		if newx < 0 || newx >= len(world[0]) || newy < 0 || newy >= len(world) {
			break
		}

		if seen.SeenDirection(newx, newy, direction) {
			// been here before, moving in this direction as well
			return -1
		}

		if world[newy][newx] == '#' {
			direction = direction.TurnRight()
		} else {
			locx, locy = newx, newy
			seen.Visit(locx, locy, direction)
		}
	}

	count := 0

	for y, line := range world {
		for x := range line {
			if seen.SeenAny(x, y) {
				count++
			}
		}
	}

	return count
}

func main() {
	world := [][]rune{}

	scanner := bufio.NewScanner(os.Stdin)

	var locx, locy int

	for scanner.Scan() {
		world = append(world, []rune(scanner.Text()))

		y := len(world) - 1

		for x, rune := range world[y] {
			if rune == '^' {
				locy, locx = y, x
			}
		}
	}

	part1 := CountSteps(world, locx, locy)
	part2 := 0

	for _, line := range world {
		for x, original := range line {

			if original == '^' {
				continue
			}

			line[x] = '#'

			if CountSteps(world, locx, locy) == -1 {
				part2++
			}

			line[x] = original
		}
	}

	println(part1, part2)
}
