package main

import (
	"io"
	"os"
	"slices"
	"strings"
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

const (
	Box      rune = 'O'
	BoxLeft  rune = '['
	BoxRight rune = ']'
	Wall     rune = '#'
	Space    rune = '.'
	Player   rune = '@'
)

func DirectionFromRune(r rune) Direction {
	switch r {
	case '^':
		return Up
	case '>':
		return Right
	case '<':
		return Left
	case 'v':
		return Down
	default:
		panic("bad direction rune: " + string(r))
	}
}

func (d Direction) String() string {
	switch d {
	case Up:
		return "up"
	case Right:
		return "right"
	case Down:
		return "down"
	case Left:
		return "left"
	default:
		panic("lmao")
	}
}

func (d Direction) DxDy() (int, int) {
	switch d {
	case Up:
		return 0, -1
	case Right:
		return 1, 0
	case Down:
		return 0, 1
	case Left:
		return -1, 0
	default:
		panic("bad direction")
	}
}

func (d Direction) Add(x, y int) (int, int) {
	dx, dy := d.DxDy()

	return x + dx, y + dy
}

type WallOrGap struct {
	block rune
	x, y  int
}

func WhatIsAtTheEndOfTheBoxes(grid [][]rune, x, y int, dir Direction) WallOrGap {
	// we're surrounded with #, so no need to check for out of bounds in this loop
	for x, y = dir.Add(x, y); grid[y][x] == 'O'; x, y = dir.Add(x, y) {
	}

	return WallOrGap{
		block: grid[y][x],
		x:     x,
		y:     y,
	}
}

func IsBoxPushable(grid [][]rune, x, y int, dir Direction) bool {
	_, dy := dir.DxDy()
	vertical := dy != 0

	if !vertical {
		nx, ny := dir.Add(x, y)
		for grid[ny][nx] == BoxLeft || grid[ny][nx] == BoxRight {
			x, y = nx, ny
		}

		return grid[ny][nx] == Space
	}

	panic("TODO")
}

func SumCoordsBoxes(grid [][]rune) int {
	sum := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == Box {
				sum += y*100 + x
			}
		}
	}

	return sum
}

func Enlarge(grid [][]rune) [][]rune {
	larger := make([][]rune, len(grid))

	for _, row := range grid {
		strrow := string(row)

		strrow = strings.ReplaceAll(strrow, "#", "##")
		strrow = strings.ReplaceAll(strrow, "O", "[]")
		strrow = strings.ReplaceAll(strrow, ".", "..")
		strrow = strings.ReplaceAll(strrow, ".", "@.")

		larger = append(larger, []rune(strrow))
	}

	return larger
}

func FilterPlayer(grid [][]rune) ([][]rune, int, int) {
	g := slices.Clone(grid)

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == Player {
				g[y][x] = Space
				return g, x, y
			}
		}
	}

	panic("no player found")
}

func main() {
	input, _ := io.ReadAll(os.Stdin)
	blocks := strings.Split(string(input), "\n\n")

	world := [][]rune{}
	directions := []Direction{}

	for _, line := range strings.Split(blocks[0], "\n") {
		world = append(world, []rune(line))
	}

	for _, dir := range blocks[1] {
		if dir == '\n' {
			continue
		}
		directions = append(directions, DirectionFromRune(dir))
	}

	println(Part1(world, directions))
}

func Part1(grid [][]rune, directions []Direction) int {
	world, x, y := FilterPlayer(grid)

	for _, direction := range directions {

		nx, ny := direction.Add(x, y)

		if world[ny][nx] == Box {
			what := WhatIsAtTheEndOfTheBoxes(world, x, y, direction)

			if what.block == Space {
				// we can move, so put the box that we want to move on at the end of the stack
				world[what.y][what.x], world[ny][nx] = world[ny][nx], world[what.y][what.x]
				x, y = nx, ny
			}
		} else if world[ny][nx] == Space {
			x, y = nx, ny
		}
	}

	return SumCoordsBoxes(world)
}
