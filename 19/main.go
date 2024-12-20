package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

func main() {
	input, _ := io.ReadAll(os.Stdin)
	split := strings.Split(string(input), "\n\n")
	towels := strings.Split(split[0], ", ")
	combinations := strings.Split(split[1], "\n")

	combinations = combinations

	slices.Sort(towels)

	fmt.Printf("%d %+v\n", len(towels), towels)
}

func Possible(towels []string, combination string) {
	used := map[string]bool{}
	longest := 0

	for _, towel := range towels {
		used[towel] = false
		longest = max(longest, len(towel))
	}

}
