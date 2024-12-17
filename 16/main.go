package main

import (
	"container/heap"
	"io"
	"os"
	"slices"
	"strings"
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type Location struct {
	x, y int
}

func (d Direction) Add(x, y int) (int, int) {
	switch d {
	case North:
		return x, y - 1
	case East:
		return x + 1, y
	case South:
		return x, y + 1
	case West:
		return x - 1, y
	default:
		panic("bad direction")
	}
}

func (d Direction) RotateLeft() Direction {
	return []Direction{West, North, East, South}[int(d)]
}

func (d Direction) RotateRight() Direction {
	return []Direction{East, South, West, North}[int(d)]
}

type Position struct {
	Location
	direction Direction
}

func (p Position) Step() Position {
	x, y := p.direction.Add(p.x, p.y)
	return Position{
		Location:  Location{x, y},
		direction: p.direction,
	}
}

func (p Position) TurnLeft() Position {
	return Position{
		Location:  p.Location,
		direction: p.direction.RotateLeft(),
	}
}

func (p Position) TurnRight() Position {
	return Position{
		Location:  p.Location,
		direction: p.direction.RotateRight(),
	}
}

type Journey struct {
	Position
	cost int
}

type PriorityQueue []Journey

func (q PriorityQueue) Len() int {
	return len(q)
}

func (q PriorityQueue) Less(i, j int) bool {
	return q[i].cost < q[j].cost
}

func (q *PriorityQueue) Push(x any) {
	*q = append(*q, x.(Journey))
}

func (q *PriorityQueue) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}

func (q PriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// compile time assertion that I implement the required methods
var _ heap.Interface = &PriorityQueue{}

func Index[T any](grid [][]T, loc Location) T {
	return grid[loc.y][loc.x]
}

func main() {
	input, _ := io.ReadAll(os.Stdin)
	maze := [][]rune{}

	var start, end Location

	for y, line := range strings.Split(string(input), "\n") {
		maze = append(maze, []rune(line))

		if x := strings.Index(line, "S"); x != -1 {
			start = Location{x: x, y: y}
		}

		if x := strings.Index(line, "E"); x != -1 {
			end = Location{x: x, y: y}
		}
	}

	cost, prev := WalkWithExpensiveCorners(maze, start)
	println(slices.Min(IndexAllDirections(cost, end)))

	// backtrack from end, save all visited locations
	visited := map[Location]bool{}

	for _, positions := range IndexAllDirections(prev, end) {
		for _, position := range positions {
			Backtrack(&visited, prev, position)
		}
	}

	println(len(visited))
}

func IndexAllDirections[T any](container map[Position]T, location Location) []T {
	result := []T{}
	for _, d := range []Direction{North, East, South, West} {
		idx := Position{
			Location:  location,
			direction: d,
		}

		if thingatpos, ok := container[idx]; ok {
			result = append(result, thingatpos)
		}
	}

	return result
}

func Backtrack(visited *map[Location]bool, prev map[Position][]Position, pos Position) {
	for _, position := range prev[pos] {
		if pos == position {
			panic("prev pos equals current pos")
		}
		(*visited)[position.Location] = true

		Backtrack(visited, prev, position)
	}
}

func WalkWithExpensiveCorners(maze [][]rune, start Location) (map[Position]int, map[Position][]Position) {
	moves := &PriorityQueue{Journey{
		Position: Position{
			Location:  start,
			direction: East,
		},
		cost: 0,
	}}

	costs := map[Position]int{
		{start, East}: 0,
	}
	prev := map[Position][]Position{}

	for moves.Len() > 0 {
		current := heap.Pop(moves).(Journey)

		neighbours := []Journey{
			{
				Position: current.Position.TurnLeft(),
				cost:     current.cost + 1000,
			},
			{
				Position: current.Position.TurnRight(),
				cost:     current.cost + 1000,
			},
		}

		if Index(maze, current.Step().Location) == '.' || Index(maze, current.Step().Location) == 'E' {
			neighbours = append(neighbours, Journey{
				Position: current.Position.Step(),
				cost:     current.cost + 1,
			})
		}

		for _, neighbour := range neighbours {
			if oldcost, ok := costs[neighbour.Position]; !ok || neighbour.cost < oldcost {
				heap.Push(moves, neighbour)
				costs[neighbour.Position] = neighbour.cost
				prev[neighbour.Position] = []Position{current.Position}
			} else if neighbour.cost == oldcost {
				prev[neighbour.Position] = append(prev[neighbour.Position], current.Position)
			}
		}
	}

	return costs, prev
}
