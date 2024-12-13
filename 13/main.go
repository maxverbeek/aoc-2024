package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Tuple2 struct {
	x, y int
}

func main() {
	allinput, _ := io.ReadAll(os.Stdin)
	blocks := strings.Split(string(allinput), "\n\n")

	part1, part2 := 0, 0

	for _, block := range blocks {
		var button1, button2, dest Tuple2

		fmt.Sscanf(block, "Button A: X+%d, Y+%d\nButton B: X+%d, Y+%d\nPrize: X=%d, Y=%d\n", &button1.x, &button1.y, &button2.x, &button2.y, &dest.x, &dest.y)

		b1, b2 := Solve(button1, button2, dest)

		dest2 := Tuple2{x: dest.x + 10000000000000, y: dest.y + 10000000000000}

		b3, b4 := Solve(button1, button2, dest2)

		part1 += b1*3 + b2
		part2 += b3*3 + b4
	}

	println(part1, part2)
}

func Solve(button1, button2, dest Tuple2) (int, int) {
	// aA + bB = dest
	// a x1 + b x2 = Dx
	// a y1 + b y2 = Dy

	// determinant = x1 * y2 - x2 * y1
	d := button1.x*button2.y - button2.x*button1.y

	if d == 0 {
		return 0, 0
	}

	t1d := dest.x*button2.y - dest.y*button2.x
	t2d := dest.y*button1.x - dest.x*button1.y

	// ensure integer solutions
	if t1d%d != 0 || t2d%d != 0 {
		return 0, 0
	}

	return t1d / d, t2d / d
}
