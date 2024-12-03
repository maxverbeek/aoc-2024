package main

import (
	. "fmt"
	. "io"
	"os"
	. "regexp"
)

func mul(line string) (sum int) {
	// inline the regex to save a variable assignment, which saves tokens
	for _, match := range MustCompile(`mul\((\d+),(\d+)\)`).FindAllString(line, -1) {
		var num1, num2 int
		Sscanf(match, "mul(%d,%d)", &num1, &num2)
		sum += num1 * num2
	}

	return
}

func main() {
	input, _ := ReadAll(os.Stdin)
	lines := string(input)

	part2 := 0

	for _, match := range MustCompile(`(?s)(?:^|do\(\))(.*?)(?:don't\(\)|$)`).FindAllString(lines, -1) {
		part2 += mul(match)
	}

	println(mul(lines), part2)
}
