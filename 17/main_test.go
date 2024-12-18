package main

import (
	"slices"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	input := "Register A: 729\nRegister B: 0\nRegister C: 0\n\nProgram: 0,1,5,4,3,0\n"

	state := Parse(strings.NewReader(input))

	if state.A != 729 || state.B != 0 || state.C != 0 || !slices.Equal(state.program, []int64{0, 1, 5, 4, 3, 0}) {
		t.Error("wrong state parsed")
	}
}

func TestExample1(t *testing.T) {
	state := State{
		A:       729,
		B:       0,
		C:       0,
		program: []int64{0, 1, 5, 4, 3, 0},
		output:  []int64{},
	}

	RunToHalt(&state)

	expected := "4,6,3,5,6,3,5,2,1,0"

	if JoinCommas(state.output) != expected {
		t.Errorf("wrong output: produced %s expected %s", JoinCommas(state.output), expected)
	}
}

func TestExample2(t *testing.T) {
	state := State{
		A:       0,
		B:       0,
		C:       9,
		program: []int64{2, 6},
		output:  []int64{},
	}

	RunToHalt(&state)

	if state.B != 1 {
		t.Errorf("Expected register B to be 1, register B is %d", state.B)
	}
}

func TestExample3(t *testing.T) {
	state := State{
		A:       10,
		program: []int64{5, 0, 5, 1, 5, 4},
		output:  []int64{},
	}

	RunToHalt(&state)

	if !slices.Equal(state.output, []int64{0, 1, 2}) {
		t.Fail()
	}
}

func TestExample4(t *testing.T) {
	state := State{
		A:       2024,
		program: []int64{0, 1, 5, 4, 3, 0},
		output:  []int64{},
	}

	RunToHalt(&state)

	expectedOutput := []int64{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0}
	if !slices.Equal(state.output, expectedOutput) {
		t.Errorf("Output mismatch. Got: %v, Expected: %v", state.output, expectedOutput)
	}
}

func TestExample5(t *testing.T) {
	state := State{
		B:       29,
		program: []int64{1, 7},
		output:  []int64{},
	}

	RunToHalt(&state)

	expectedB := int64(26)
	if state.B != expectedB {
		t.Errorf("Register B mismatch. Got: %d, Expected: %d", state.B, expectedB)
	}
}

func TestExample6(t *testing.T) {
	state := State{
		B:       2024,
		C:       43690,
		program: []int64{4, 0},
		output:  []int64{},
	}

	RunToHalt(&state)

	expectedB := int64(44354)
	if state.B != expectedB {
		t.Errorf("Register B mismatch. Got: %d, Expected: %d", state.B, expectedB)
	}
}
