package main

import (
	"bufio"
	"fmt"
	"os"
)

func Mix(num, val int) int {
	return num ^ val
}

func Prune(num int) int {
	return num % 16777216
}

func NextSecret(num int) int {
	num = Prune(Mix(num, num*64))
	num = Prune(Mix(num, num/32))
	num = Prune(Mix(num, num*2048))

	return Prune(num)
}

func AfterNSecret(num, n int) int {
	for i := 0; i < n; i++ {
		num = NextSecret(num)
	}
	return num
}

type Sequence struct {
	ab, bc, cd, de int
}

func BuildSequences(init, n int) map[Sequence]int {
	a := init
	b := NextSecret(a)
	c := NextSecret(b)
	d := NextSecret(c)
	e := NextSecret(d)
	results := map[Sequence]int{}

	for i := 4; i < n; i++ {
		sequence := Sequence{
			b%10 - a%10,
			c%10 - b%10,
			d%10 - c%10,
			e%10 - d%10,
		}

		// we only care about the first occurence of the sequence, so ignore
		// any future ones (which will overwrite the value of previously found
		// sequences)
		if _, exists := results[sequence]; !exists {
			results[sequence] = e % 10
		}

		a, b, c, d, e = b, c, d, e, NextSecret(e)
	}

	return results
}

func SumFromKey[K comparable](iterablemaps []map[K]int, key K) int {
	sum := 0
	for _, m := range iterablemaps {
		sum += m[key]
	}

	return sum
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	part1 := 0

	sequences := []map[Sequence]int{}

	allsequences := map[Sequence]struct{}{}

	for scanner.Scan() {
		var initial int

		fmt.Sscanf(scanner.Text(), "%d", &initial)

		part1 += AfterNSecret(initial, 2000)

		mysequences := BuildSequences(initial, 2000)
		sequences = append(sequences, mysequences)

		for seq := range mysequences {
			allsequences[seq] = struct{}{}
		}
	}

	println(part1)

	maxgains := -9999

	for seq := range allsequences {
		mygains := SumFromKey(sequences, seq)

		if mygains > maxgains {
			maxgains = mygains
		}
	}

	println(maxgains)
}
