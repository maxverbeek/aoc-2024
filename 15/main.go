package main

import (
	"io"
	"os"
	"strings"
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Block rune

const (
	Box      Block = 'O'
	BoxLeft  Block = '['
	BoxRight Block = ']'
	Wall     Block = '#'
	Space    Block = '.'
	Player   Block = '@'
)

func (b Block) IsBox() bool {
	return b == Box || b == BoxLeft || b == BoxRight
}

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
	block Block
	x, y  int
}

func WhatIsAtTheEndOfTheBoxes(grid [][]Block, x, y int, dir Direction) WallOrGap {
	// we're surrounded with #, so no need to check for out of bounds in this loop
	for x, y = dir.Add(x, y); grid[y][x] == Box; x, y = dir.Add(x, y) {
	}

	return WallOrGap{
		block: grid[y][x],
		x:     x,
		y:     y,
	}
}

func CanPushBoxes(grid [][]Block, x, y int, dir Direction) bool {
	nx, ny := dir.Add(x, y)
	vertical := y != ny

	if grid[ny][nx] == Space {
		return true
	}

	if grid[ny][nx] == Wall {
		return false
	}

	if !grid[ny][nx].IsBox() {
		panic("somehow ended up not at a box")
	}

	if !vertical {
		return CanPushBoxes(grid, nx, ny, dir)
	}

	// we are pushing boxes vertically, so assign the left half to l and the right to r
	var lx, ly, rx, ry int

	if grid[ny][nx] == BoxLeft {
		lx, ly = nx, ny
		rx, ry = nx+1, ny
	} else {
		lx, ly = nx-1, ny
		rx, ry = nx, ny
	}

	// with overlapping boxes, this has the potential to make lots of duplicate calls
	// FIXME: cache this
	return CanPushBoxes(grid, lx, ly, dir) && CanPushBoxes(grid, rx, ry, dir)
}

func MoveBoxes(grid [][]Block, x, y int, dir Direction) {
	nx, ny := dir.Add(x, y)
	vertical := y != ny

	if grid[ny][nx] == Space {
		grid[ny][nx], grid[y][x] = grid[y][x], grid[ny][nx]
		return
	}

	if !vertical && grid[ny][nx].IsBox() {
		// first move the next boxes
		MoveBoxes(grid, nx, ny, dir)
		// then move the current box into the newly created space
		grid[ny][nx], grid[y][x] = grid[y][x], grid[ny][nx]

		return
	}

	var lx, ly, rx, ry int
	if grid[ny][nx] == BoxLeft {
		lx, ly = nx, ny
		rx, ry = nx+1, ny
	} else {
		lx, ly = nx-1, ny
		rx, ry = nx, ny
	}

	// move the boxes up ahead
	MoveBoxes(grid, lx, ly, dir)
	MoveBoxes(grid, rx, ry, dir)

	// move the current box into its newly created space
	grid[ny][nx], grid[y][x] = grid[y][x], grid[ny][nx]
}

func SumCoordsBoxes(grid [][]Block) int {
	sum := 0
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == Box || grid[y][x] == BoxLeft {
				sum += y*100 + x
			}
		}
	}

	return sum
}

func Enlarge(input string) string {
	input = strings.ReplaceAll(input, "#", "##")
	input = strings.ReplaceAll(input, "O", "[]")
	input = strings.ReplaceAll(input, ".", "..")
	input = strings.ReplaceAll(input, "@", "@.")

	return input
}

func FilterPlayer(input string) ([][]Block, int, int) {
	world := [][]Block{}

	var x, y int

	for i, line := range strings.Split(input, "\n") {
		world = append(world, []Block(line))

		if j := strings.Index(line, "@"); j >= 0 {
			y = i
			x = j
			world[y][x] = Space
		}
	}

	if x == 0 || y == 0 {
		panic("no player found")
	}

	return world, x, y
}

func Draw(grid [][]Block) {
	for _, line := range grid {
		println(string(line))
	}
}

func main() {
	input, _ := io.ReadAll(os.Stdin)
	blocks := strings.Split(string(input), "\n\n")

	world, x, y := FilterPlayer(blocks[0])
	directions := []Direction{}

	for _, dir := range blocks[1] {
		if dir == '\n' {
			continue
		}
		directions = append(directions, DirectionFromRune(dir))
	}

	println(Part1(world, x, y, directions))

	world, x, y = FilterPlayer(Enlarge(blocks[0]))

	println(Part2(world, x, y, directions))
}

func Part1(world [][]Block, x, y int, directions []Direction) int {
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

func Part2(world [][]Block, x, y int, directions []Direction) int {
	for _, direction := range directions {
		nx, ny := direction.Add(x, y)

		if world[ny][nx].IsBox() && CanPushBoxes(world, x, y, direction) {
			MoveBoxes(world, x, y, direction)
			x, y = nx, ny
		} else if world[ny][nx] == Space {
			x, y = nx, ny
		}
	}

	return SumCoordsBoxes(world)
}
