package main

import (
	"io"
	"os"
	"strings"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)
	split := strings.Split(string(input), "\n\n")
	towels := strings.Split(split[0], ", ")
	combinations := strings.Split(strings.TrimSpace(split[1]), "\n")

	usedtowels := map[string]bool{}
	previouscombinations := map[string]int{}

	for _, towel := range towels {
		usedtowels[towel] = false
	}

	part1 := 0
	part2 := 0

	for _, combination := range combinations {
		if p := Possibilities(usedtowels, combination, &previouscombinations); p > 0 {
			part1++
			part2 += p
		}
	}

	println(part1, part2)
}

func Possibilities(towels map[string]bool, combination string, combinations *map[string]int) int {
	if len(combination) == 0 {
		return 1
	}

	if c, ok := (*combinations)[combination]; ok {
		return c
	}

	p := 0

	defer func() {
		(*combinations)[combination] = p
	}()

	for towel := range towels {
		if strings.HasPrefix(combination, towel) {
			p += Possibilities(towels, combination[len(towel):], combinations)
		}
	}

	return p
}
