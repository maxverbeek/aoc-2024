package main

import (
	"bufio"
	"fmt"
	"os"
)

type Tuple = struct {
	x, y int
}

var (
	numpad = map[rune]Tuple{
		'7': {0, 0},
		'8': {1, 0},
		'9': {2, 0},
		'4': {0, 1},
		'5': {1, 1},
		'6': {2, 1},
		'1': {0, 2},
		'2': {1, 2},
		'3': {2, 2},
		' ': {0, 3},
		'0': {1, 3},
		'A': {2, 3},
	}

	dirpad = map[rune]Tuple{
		' ': {0, 0},
		'^': {1, 0},
		'A': {2, 0},
		'<': {0, 1},
		'v': {1, 1},
		'>': {2, 1},
	}
)

func CombineGreedy(keypad map[rune]Tuple, combination []rune) []rune {
	result := []rune{}
	from := 'A'

	for i := range combination {
		to := combination[i]

		fx, fy := keypad[from].x, keypad[from].y
		tx, ty := keypad[to].x, keypad[to].y

		hfirst, vfirst := Tuple{x: tx, y: fy}, Tuple{x: fx, y: ty}

		// the ideal order is left, down, right up, ordered by decreasing distance from the A button
		// but check if we can go left first due to the empty square
		// since the empty square is in top-left, or bottom left for the numpad,
		// we never have to worry about crossing it when going right
		if tx < fx && keypad[' '] != hfirst {
			// want to go left and can go left
			for tx < fx {
				result = append(result, '<')
				fx--
			}
		}

		if ty > fy && keypad[' '] != vfirst {
			// want to go down and can go down
			for ty > fy {
				result = append(result, 'v')
				fy++
			}
		}

		if ty < fy && keypad[' '] != vfirst {
			// want to go up and can go up
			for ty < fy {
				result = append(result, '^')
				fy--
			}
		}

		for tx > fx {
			// going right is always safe
			result = append(result, '>')
			fx++
		}

		// if we couldn't go left at first, we certainly can now because both
		// vertical directions are possible and have been executed if the
		// horizontal direction was problematic.
		for tx < fx {
			result = append(result, '<')
			fx--
		}

		// same for going down: left and right have now definitely been covered
		for ty > fy {
			result = append(result, 'v')
			fy++
		}

		// going up as well
		for ty < fy {
			result = append(result, '^')
			fy--
		}

		result = append(result, 'A')
		from = to
	}

	return result
}

func EncodePadsGreedy(pads []map[rune]Tuple, combination []rune) []rune {
	for i, pad := range pads {
		fmt.Printf("encoding pad %d with length %d\n", i, len(combination))
		combination = CombineGreedy(pad, combination)
	}

	return combination
}

func main() {
	part1, part2 := 0, 0
	scanner := bufio.NewScanner(os.Stdin)

	part1maps := []map[rune]Tuple{numpad, dirpad, dirpad}
	part2maps := []map[rune]Tuple{numpad}

	for i := 0; i < 26; i++ {
		part2maps = append(part2maps, dirpad)
	}

	for scanner.Scan() {
		var num int
		fmt.Sscanf(scanner.Text(), "%dA", &num)

		part1input := EncodePadsGreedy(part1maps, []rune(scanner.Text()))
		part2input := EncodePadsGreedy(part2maps, []rune(scanner.Text()))

		part1 += len(part1input) * num
		part2 += len(part2input) * num
	}

	println(part1, part2)
}
