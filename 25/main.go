package main

import (
	"io"
	"os"
	"strings"
)

type Pattern struct {
	height  int
	decoded []int
}

func main() {
	input, _ := io.ReadAll(os.Stdin)
	keysorlocks := strings.Split(string(input), "\n\n")

	keys, locks := []Pattern{}, []Pattern{}

	for _, keyorlock := range keysorlocks {
		if IsLock(keyorlock) {
			locks = append(locks, Decode(keyorlock))
		} else {
			keys = append(keys, Decode(keyorlock))
		}
	}

	count := 0
	for _, lock := range locks {
		for _, key := range keys {
			overlapping := false
			if key.height != lock.height {
				panic("key and lock mismatching height")
			}

			for pin := range key.decoded {
				if key.decoded[pin]+lock.decoded[pin] >= key.height {
					overlapping = true
				}
			}

			if !overlapping {
				count++
			}
		}
	}

	println(count)
}

func IsLock(pattern string) bool {
	lines := strings.Split(pattern, "\n")

	isLock := true
	isKey := true

	for _, c := range []rune(lines[0]) {
		isLock = isLock && c == '#'
	}

	for _, c := range []rune(lines[len(lines)-1]) {
		isKey = isKey && c == '#'
	}

	return isLock
}

func Decode(pattern string) Pattern {
	result := []int{}

	lines := strings.Split(strings.TrimSpace(pattern), "\n")

	width, height := len(lines[0]), len(lines)

	for x := 0; x < width; x++ {
		sum := 0
		for y := 0; y < height; y++ {
			if lines[y][x] == '#' {
				sum++
			}
		}

		result = append(result, sum-1)
	}

	return Pattern{
		decoded: result,
		height:  height - 1,
	}
}
