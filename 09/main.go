package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

type BlockType int

const (
	File BlockType = iota
	Space
)

type Block struct {
	blocktype BlockType
	id        int
	size      int
}

func Print(disk *list.List) {
	for block := disk.Front(); block != nil; block = block.Next() {
		b := block.Value.(Block)

		for i := 0; i < b.size; i++ {
			if b.blocktype == File {
				fmt.Printf("%d", b.id)
			} else {
				fmt.Printf(".")
			}
		}
	}

	fmt.Printf("\n")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()

	disk := list.New()
	disk2 := list.New()

	for i, char := range line {
		block := Block{
			size: int(char - '0'),
		}

		if i%2 == 0 {
			// file
			block.id = i / 2
			block.blocktype = File
		} else {
			block.blocktype = Space
		}

		disk.PushBack(block)
		disk2.PushBack(block)
	}

	println(Part1(disk))
	println(Part2(disk2))
}

func Part1(disk *list.List) int {
	start, end := disk.Front(), disk.Back()

	for start != nil && end != nil && start != end {
		// Print(disk)
		startblock, endblock := start.Value.(Block), end.Value.(Block)

		if startblock.blocktype != Space {
			start = start.Next()
			continue
		}

		if endblock.blocktype != File {
			end = end.Prev()
			continue
		}

		// invariant: start is space, end is a file
		if startblock.size < endblock.size {
			// end is too large, so split it in two
			splitblock := Block{
				id:        endblock.id,
				blocktype: endblock.blocktype,
				size:      endblock.size - startblock.size,
			}
			disk.InsertBefore(splitblock, end)
			endblock.size = startblock.size
		} else if startblock.size > endblock.size {
			// split up start because it's too big
			splitblock := Block{
				id:        startblock.id,
				blocktype: startblock.blocktype,
				size:      startblock.size - endblock.size,
			}
			disk.InsertAfter(splitblock, start)
			startblock.size = endblock.size
		}

		// size matches exactly
		start.Value = endblock
		end.Value = Block{blocktype: Space}
	}

	return Count(disk)
}

func Part2(disk *list.List) int {
	for end := disk.Back(); end != disk.Front(); end = end.Prev() {
		if end.Value.(Block).blocktype != File {
			continue
		}

		for start := disk.Front(); start != end; start = start.Next() {
			startblock, endblock := start.Value.(Block), end.Value.(Block)

			if startblock.blocktype != Space {
				continue
			}

			if startblock.size >= endblock.size {
				// end fits in start

				if startblock.size > endblock.size {
					// split start so that it fits end exactly
					split := Block{
						id:        startblock.id,
						blocktype: startblock.blocktype,
						size:      startblock.size - endblock.size,
					}

					disk.InsertAfter(split, start)
					startblock.size -= split.size
					start.Value = startblock
				}

				// move end to start
				start.Value, end.Value = end.Value, start.Value

				break
			}
		}
	}

	return Count(disk)
}

func Count(disk *list.List) int {
	var position, count int

	for block := disk.Front(); block != nil; block = block.Next() {
		b := block.Value.(Block)

		for i := 0; i < b.size; i++ {
			count += position * b.id
			position++
		}
	}

	return count
}
