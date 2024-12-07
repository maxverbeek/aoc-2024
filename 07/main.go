package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ParseNumbers(numstring string) []int64 {
	unparsednumbers := strings.Split(numstring, " ")
	numbers := make([]int64, len(unparsednumbers))

	for i, text := range unparsednumbers {
		fmt.Sscanf(text, "%d", &numbers[i])
	}

	return numbers
}

func MatchMulOrAdd(current, target int64, operands []int64) bool {
	if len(operands) == 0 {
		return current == target
	}

	next := operands[0]

	return MatchMulOrAdd(current*next, target, operands[1:]) || MatchMulOrAdd(current+next, target, operands[1:])
}

func MatchMulAddOrConcat(current, target int64, operands []int64) bool {
	if len(operands) == 0 {
		return current == target
	}

	next := operands[0]

	var concatted int64
	fmt.Sscanf(fmt.Sprintf("%d%d", current, next), "%d", &concatted)

	return MatchMulAddOrConcat(current*next, target, operands[1:]) ||
		MatchMulAddOrConcat(current+next, target, operands[1:]) ||
		MatchMulAddOrConcat(concatted, target, operands[1:])
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var part1, part2 int64

	for scanner.Scan() {
		segments := strings.Split(scanner.Text(), ": ")
		testvalues, operands := ParseNumbers(segments[0]), ParseNumbers(segments[1])
		testvalue := testvalues[0]

		if MatchMulOrAdd(operands[0], testvalue, operands[1:]) {
			part1 += testvalue
		}

		if MatchMulAddOrConcat(operands[0], testvalue, operands[1:]) {
			part2 += testvalue
		}
	}

	println(part1, part2)
}
