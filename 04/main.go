package main

import (
	"bufio"
	"fmt"
	"os"
)

var input []string
var width, height int

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	width, height = len(input[0]), len(input)
	var part1, part2 int

	for row, line := range input {
		for col := range line {
			var deltas = [...]int{-1, 0, 1}
			for _, dx := range deltas {
				for _, dy := range deltas {
					part1 += Search1(row, col, dx, dy)
				}
			}
			part2 += Search2(row, col)
		}
	}

	println(part1, part2)
}

func Search1(row, col, dx, dy int) int {
	for i, r := range "XMAS" {
		x := col + dx*i
		y := row + dy*i

		if x < 0 || x >= width || y < 0 || y >= height || input[y][x] != byte(r) {
			return 0
		}
	}

	return 1
}

// 1899 is too low
func Search2(row, col int) int {
	if input[row][col] != 'A' {
		return 0
	}

	if row < 1 || row >= height-1 || col < 1 || col >= width-1 {
		return 0
	}

	d1 := fmt.Sprintf("%c%c", input[row-1][col-1], input[row+1][col+1])
	d2 := fmt.Sprintf("%c%c", input[row+1][col-1], input[row-1][col+1])

	if (d1 != "MS" && d1 != "SM") || (d2 != "MS" && d2 != "SM") {
		return 0
	}

	return 1
}
