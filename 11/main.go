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

	stones := make([]int64, len(words))

	for i, word := range words {
		stone, _ := strconv.ParseInt(word, 10, 64)
		stones[i] = stone
	}

	for blink := 0; blink < 25; blink++ {
		stones = NewStones(stones)
	}

	println(len(stones))
}

func Log10(num int64) int64 {
	log := int64(1)

	for num >= 10 {
		log++
		num /= 10
	}

	return log
}

func Pow10(pow int64) int64 {
	num := int64(1)
	for pow > 0 {
		num *= 10
		pow--
	}

	return num
}

func NewStones(oldstones []int64) []int64 {
	newstones := []int64{}

	for _, stone := range oldstones {
		digitcount := Log10(stone)
		if stone == 0 {
			newstones = append(newstones, 1)
		} else if digitcount%2 == 0 {
			halfpow := Pow10(digitcount / 2)
			newstones = append(newstones, stone/halfpow, stone%halfpow)
		} else {
			newstones = append(newstones, stone*2024)
		}
	}

	return newstones
}
