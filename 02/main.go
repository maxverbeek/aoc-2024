package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// func main() {
// 	var num int
// 	var whitespace rune
// 	_, err := fmt.Scanf("%d%c", &num, &whitespace)
// 	for err == nil && whitespace != '\n' {
// 		println(num)
// 		_, err = fmt.Scanf("%d%c", &num, &whitespace)
// 	}
// }

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	valid := 0

outter:
	for scanner.Scan() {
		levels := strings.Fields(scanner.Text())
		println(scanner.Text())
		ascending := true
		for i := 1; i < len(levels); i++ {
			level1, _ := strconv.Atoi(levels[i-1])
			level2, _ := strconv.Atoi(levels[i])
			diff := max(level1-level2, level2-level1)
			if diff < 1 || diff > 3 {
				fmt.Printf("rejected by %f %f difference too big\n", level1, level2)
				continue outter
			}

			if i == 1 {
				ascending = level2 > level1
			} else if ascending != (level2 > level1) {
				fmt.Printf("rejected by %f %f (%b)\n", level1, level2, ascending)
				continue outter
			}
		}
		valid++
	}

	println(valid)
}
