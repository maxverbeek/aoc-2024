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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	line := scanner.Text()

	disk := list.New()

	for i, char := range line {
		block := Block{
			size: int(char - '0'),
		}

		if i%2 == 0 {
			// file
			block.id = i
			block.blocktype = File
		} else {
			block.blocktype = Space
		}

		disk.PushBack(block)
	}

	for block := disk.Front(); block != nil; block = block.Next() {
		b := block.Value.(Block)

		if b.blocktype == File {
			fmt.Printf("%d", b.id)
		} else {
			fmt.Printf(".")
		}
	}

	start, end := disk.Front(), disk.Back()

	for start != nil && end != nil && start != end {
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
			// split up end because it's too big
			split := Block{
				id:        startblock.id,
				blocktype: startblock.blocktype,
				size:      endblock.size - startblock.size,
			}
			disk.InsertAfter(split, start)
			startblock.size = endblock.size
			start.Value = startblock
		} else if startblock.size > endblock.size {
			// split up start because it's too big
			split := Block{
				id:        endblock.id,
				blocktype: endblock.blocktype,
				size:      startblock.size - endblock.size,
			}
			fmt.Printf("Splitting start block: %v\n", end)
			disk.InsertBefore(split, end)
			endblock.size = startblock.size
			end.Value = endblock
		}

		// size matches exactly
		start.Value = endblock
		end.Value = Block{blocktype: Space}
	}

	part1 := 0

	position := -1

	for block := disk.Front(); block != nil; block = block.Next() {
		position++
		b := block.Value.(Block)

		if b.blocktype == File {
			fmt.Printf("%d", b.id)
		} else {
			fmt.Printf(".")
		}

		if b.blocktype == Space {
			continue
		}

		part1 += position * b.size
	}

	println(part1)
}
