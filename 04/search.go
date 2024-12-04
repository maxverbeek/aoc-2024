package main

import . "strings"

func Search(kerneltemplates string) (count int) {
	// split by 2 new lines to obtain a block of an actual kernel
	for _, kernellines := range Split(kerneltemplates, "\n\n") {

		// further split this block into different lines so we can iterate over them
		kernel := Split(kernellines, "\n")

		// these loops are super nasty, abusing the range syntax to iterate
		// over the right amount of elements just to get an index that
		// corresponds to the size of the input minus the size of the kernel.
		// this inputrow stuff is not actually really used, these loops are
		// actually just this:
		//
		// for row := 0; row <= len(kernel) + len(input); row++
		//
		// for col := 0; col <= len(kernel[0]) + len(input[0]); col++
		//
		// but that is slightly more tokens
		for row, inputrow := range input[len(kernel)-1:] {
		columns:
			for col := range inputrow[len(kernel[0])-1:] {
				// here we try to fit a kernel at the current position
				for y, kernely := range kernel {
					for x, kernelrune := range kernely {
						if kernelrune != ' ' && input[row+y][col+x] != byte(kernelrune) {
							// if we find a mistake, continue with the next
							// column which will try fitting the current
							// kernel somewhere else. if this doesn't work out
							// anywhere this function moves over to the next
							// kernel
							continue columns
						}
					}
				}

				// if the loop reaches this point, the kernel fits at this
				// position because it survived all of the checks.
				count++
			}
		}
	}
	return
}
