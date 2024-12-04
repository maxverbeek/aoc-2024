package main

import (
	"bufio"
	"os"
)

var input []string

func main() {
	kernels1 := [][]string{
		{
			"XMAS",
		}, {
			"SAMX",
		},
		{
			"X",
			"M",
			"A",
			"S",
		},
		{
			"S",
			"A",
			"M",
			"X",
		},
		{
			"X   ",
			" M  ",
			"  A ",
			"   S",
		},
		{
			"   X",
			"  M ",
			" A  ",
			"S   ",
		},
		{
			"   S",
			"  A ",
			" M  ",
			"X   ",
		},
		{
			"S   ",
			" A  ",
			"  M ",
			"   X",
		},
	}

	kernels2 := [][]string{
		{
			"M S",
			" A ",
			"M S",
		},
		{
			"M M",
			" A ",
			"S S",
		},
		{
			"S M",
			" A ",
			"S M",
		},
		{
			"S S",
			" A ",
			"M M",
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	println(Search(kernels1), Search(kernels2))
}

func Search(kernels [][]string) (count int) {
	for _, kernel := range kernels {
		for row := 0; row <= len(input)-len(kernel); row++ {
		kernel:
			for col := 0; col <= len(input[0])-len(kernel[0]); col++ {
				for y := range kernel {
					for x := range kernel[y] {
						if kernel[y][x] != ' ' && input[row+y][col+x] != kernel[y][x] {
							continue kernel
						}
					}
				}
				count++
			}
		}
	}
	return
}
