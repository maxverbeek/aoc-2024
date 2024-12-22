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

type CacheKey struct {
	from rune
	rest string
}

func AddLtGt(target, coordinate int, lt, gt rune) []rune {
	result := []rune{}

	for coordinate < target {
		result = append(result, lt)
		coordinate++
	}

	for coordinate > target {
		result = append(result, gt)
		coordinate--
	}

	return result
}

func TypeCombinations(keypad map[rune]Tuple, cache map[CacheKey][][]rune, from rune, combination []rune) [][]rune {
	cachekey := CacheKey{from: from, rest: string(combination)}
	if results, ok := cache[cachekey]; ok {
		return results
	}
	prefixes := [][]rune{}

	to := combination[0]

	// horizontal means we're moving horizontally first, so x will change
	// and y will stay at the from. the corner that this makes cannot end
	// up on the empty square
	if keypad[' '] != (Tuple{x: keypad[to].x, y: keypad[from].y}) {
		prefix := []rune{}
		prefix = append(prefix, AddLtGt(keypad[to].x, keypad[from].x, '>', '<')...)
		prefix = append(prefix, AddLtGt(keypad[to].y, keypad[from].y, 'v', '^')...)
		prefix = append(prefix, 'A')

		prefixes = append(prefixes, prefix)
	}

	// !horizontal = vertical means y changes first, and x remains at the from position.
	if keypad[' '] != (Tuple{x: keypad[from].x, y: keypad[to].y}) {
		prefix := []rune{}
		prefix = append(prefix, AddLtGt(keypad[to].x, keypad[from].x, '>', '<')...)
		prefix = append(prefix, AddLtGt(keypad[to].y, keypad[from].y, 'v', '^')...)
		prefix = append(prefix, 'A')

		prefixes = append(prefixes, prefix)
	}

	if len(combination) > 1 {
		suffixes := TypeCombinations(keypad, cache, to, combination[1:])

		results := [][]rune{}

		for _, prefix := range prefixes {
			for _, suffix := range suffixes {
				results = append(results, []rune(string(prefix)+string(suffix)))
			}
		}

		cache[cachekey] = results

		fmt.Printf("made %d results\n", len(results))

		return results
	}

	return prefixes
}

func Allpads(pads []map[rune]Tuple, combination []rune) [][]rune {
	if len(pads) == 0 {
		return [][]rune{combination}
	}

	fmt.Printf("%d pads to go..\n", len(pads))

	options := TypeCombinations(pads[0], make(map[CacheKey][][]rune), 'A', combination)

	next := [][]rune{}

	minlength, mincount := len(options[0]), 0

	for i := range options {
		if len(options[i]) < minlength {
			minlength = len(options[i])
			mincount = 1
		} else if len(options[i]) == minlength {
			mincount++
		}
	}

	fmt.Printf("considering %d options on %d pads\n", mincount, len(pads))

	for _, option := range options {
		if len(option) == minlength {
			next = append(next, Allpads(pads[1:], option)...)
		}
	}

	return next
}

func RobotIndirection(combination []rune) []rune {
	combinations := Allpads([]map[rune]Tuple{numpad, dirpad, dirpad}, combination)
	minidx, minlength := len(combinations[0]), 0

	for i := range combinations {
		if len(combinations[i]) < minlength {
			minidx, minlength = i, len(combinations[i])
		}
	}

	return combinations[minidx]
}

func main() {
	part1 := 0
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var num int
		fmt.Sscanf(scanner.Text(), "%dA", &num)

		humaninput := RobotIndirection([]rune(scanner.Text()))
		part1 += len(humaninput) * num

		fmt.Printf("%s: %d*%d = %d: %s\n", scanner.Text(), len(humaninput), num, len(humaninput)*num, string(humaninput))
	}

	println(part1)
}
