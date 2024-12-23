package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestBehaviour(t *testing.T) {
	combination := []rune("1A1A")
	robot1 := CombineGreedy(numpad, 'A', combination)
	robot2 := CombineGreedy(dirpad, 'A', robot1)
	robot3 := CombineGreedy(dirpad, 'A', robot2)

	t.Logf("%s\n%s\n%s\n%s\n", string(combination), string(robot1), string(robot2), string(robot3))
}

func TestTable(t *testing.T) {
	chars := []rune("A^<v>")
	table := [5][5][]rune{}
	longest := 0

	for i := range chars {
		for j := range chars {
			table[i][j] = PathFind(dirpad, chars[i], chars[j])
			longest = max(longest, len(table[i][j]))
		}
	}

	content := bytes.NewBuffer([]byte{})

	// header
	fmt.Fprintf(content, "Table:\n")
	fmt.Fprintf(content, "  %*c %*c %*c %*c %*c\n", longest, chars[0], longest, chars[1], longest, chars[2], longest, chars[3], longest, chars[4])

	for i := range chars {
		fmt.Fprintf(content, "%c ", chars[i])
		for j := range chars {
			fmt.Fprintf(content, "%*s ", longest, string(table[i][j]))
		}
		fmt.Fprintf(content, "\n")
	}

	t.Logf(content.String())
}

func TestComputability(t *testing.T) {
	t.Skip()
	timeout := time.After(5 * time.Second)
	done := make(chan struct{})

	go func() {
		combination := []rune{'<'}

		for i := 0; i < 25; i++ {
			t.Logf("at combination %d", i)
			combination = CombineGreedy(dirpad, 'A', combination)
		}
	}()

	select {
	case <-timeout:
		t.Fatal("took too long")
	case <-done:
	}
}

// func TestCountStateTransitionsWorks(t *testing.T) {
//     ">A"
// }

func TestApproachesAreEqual(t *testing.T) {
	keypadscases := [][]map[rune]Tuple{
		{numpad, dirpad, dirpad},
		{numpad, dirpad, dirpad, dirpad, dirpad, dirpad, dirpad}, // 6 dirpads
	}

	codecases := []string{"029A", "980A", "179A", "456A", "379A"}

	for i, keypads := range keypadscases {
		for _, code := range codecases {
			naive := len(EncodePadsGreedy(keypads, []rune(code)))
			batch := BatchEncodePadsGreedy(keypads, []rune(code))

			if naive != batch {
				t.Errorf("code %s failed for keypad case %d (%d keypads): naive %d vs batch %d\n", code, i, len(keypads), naive, batch)
			}
		}
	}
}
