package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()
	words := strings.Split(line, " ")

	stones := map[int]int{}

	for _, word := range words {
		stone, _ := strconv.Atoi(word)
		stones[stone] = 1
	}

	for blink := 0; blink < 75; blink++ {
		if blink == 25 {
			println(CountHist(stones))
		}

		newstones := map[int]int{}

		for stone, count := range stones {
			for _, newstone := range EvolveStone(stone) {
				newstones[newstone] += count
			}
		}

		stones = newstones
	}

	println(CountHist(stones))
}

func CountHist(hist map[int]int) int {
	total := int(0)

	for _, count := range hist {
		total += count
	}

	return total
}

func Log10(num int) int {
	log := 1

	for num >= 10 {
		log++
		num /= 10
	}

	return log
}

func Pow10(pow int) int {
	num := 1
	for pow > 0 {
		num *= 10
		pow--
	}

	return num
}

func EvolveStone(stone int) []int {
	newstones := []int{}

	digitcount := Log10(stone)
	if stone == 0 {
		newstones = append(newstones, 1)
	} else if digitcount%2 == 0 {
		halfpow := Pow10(digitcount / 2)
		newstones = append(newstones, stone/halfpow, stone%halfpow)
	} else {
		newstones = append(newstones, stone*2024)
	}

	return newstones
}
