package main

import (
	"bufio"
	"os"
	"slices"
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

	hist := Histogram(stones)
	for blink := 0; blink < 75; blink++ {
		newstones := map[int64]int64{}

		for stone, count := range hist {
			for _, newstone := range EvolveStone(stone) {
				newstones[newstone] += count
			}
		}

		hist = newstones

		if blink == 24 {
			println(CountHist(hist))
		}
	}

	println(CountHist(hist))
}

func Largest(stones []int64) int64 {
	return slices.Max(stones)
}

func CountHist(hist map[int64]int64) int64 {
	total := int64(0)

	for _, count := range hist {
		total += count
	}

	return total
}
func Histogram(stones []int64) map[int64]int64 {
	histogram := map[int64]int64{}
	for _, stone := range stones {
		histogram[stone]++
	}

	return histogram
}

func CountUnique(stones []int64) int {
	sorted := slices.Clone(stones)
	slices.Sort(sorted)
	return len(slices.Compact(sorted))
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

func EvolveStone(stone int64) []int64 {
	newstones := []int64{}

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
