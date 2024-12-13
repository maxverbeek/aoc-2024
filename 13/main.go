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

func Gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

func Solve(button1, button2, dest Tuple2) (int, int) {

	if dest.x%Gcd(button1.x, button2.x) != 0 || dest.y%Gcd(button1.y, button2.y) != 0 {
		return 0, 0
	}

	// how many times can we fit button 2 before having to press button 1?
	times2 := min(dest.x/button2.x, dest.y/button2.y)
	times1 := 0

	for times2 >= 0 {
		times1 = min((dest.x-button2.x*times2)/button1.x, (dest.y-button2.y*times2)/button1.y)

		covered := Tuple2{
			x: button2.x*times2 + button1.x*times1,
			y: button2.y*times2 + button1.y*times1,
		}

		if covered == dest {
			return times1, times2
		}

		times2--
	}

	return 0, 0
}
