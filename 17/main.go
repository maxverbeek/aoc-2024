package main

import (
	"fmt"
	"io"
	"os"
	// "slices"
	"strconv"
	"strings"
)

type State struct {
	pc      int64
	A, B, C int64
	program []int64
	output  []int64
}

func (s *State) Combo(operand int64) int64 {
	switch operand {
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		return operand
	case 4:
		return s.A
	case 5:
		return s.B
	case 6:
		return s.C
	default:
		panic("invalid combo operand")
	}
}

func (s *State) Fetch() (bool, int64, int64) {
	if s.pc > int64(len(s.program)-2) {
		return true, 0, 0
	}

	return false, s.program[s.pc], s.program[s.pc+1]
}

func (s *State) Execute(opcode, operand int64) {
	switch opcode {
	case 0: // adv
		s.A = s.A >> s.Combo(operand)
	case 1: // bxl
		s.B = s.B ^ operand
	case 2: // bst
		s.B = s.Combo(operand) % 8
	case 3: // jnz
		if s.A != 0 {
			s.pc = operand - 2 // subtract 2 because we add 2 later
		}
	case 4: // bxc
		s.B = s.B ^ s.C
	case 5: // out
		s.output = append(s.output, s.Combo(operand)%8)
	case 6: // bdv
		s.B = s.A >> s.Combo(operand)
	case 7: // cdv
		s.C = s.A >> s.Combo(operand)
	default:
		panic(fmt.Sprintf("bad instruction: %d", opcode))
	}

	s.pc += 2
}

func ExecFast(a, b, c int64) int64 {
	out := int64(0)

	// this is my input program, translated into Go
	for a > 0 {
		// b = a & 0b111 ^ 2
		// c = a >> b
		// b = b ^ c
		// b = b ^ 3
		out = (out << 3) + ((a&0b111)^(a>>((a&0b111)^2))^3^2)&0b111
		a = a >> 3
	}

	return out
}

func (s *State) Inverted(opcode, operand int64) {
	switch opcode {
	case 0: // adv
		s.A = s.A << s.Combo(operand)
	case 1: // bxl
		s.B = s.B ^ operand
	case 2: // bst
		s.B = s.Combo(operand) % 8
	case 3: // jnz
		if len(s.output) > 0 {
			s.pc = -2
		}
	case 4: // bxc
		s.B = s.B ^ s.C
	case 5: // out
		s.A = s.A + s.output[0]
		s.output = s.output[1:]
	case 6: // bdv
		s.B = s.A << s.Combo(operand)
	case 7: // cdv
		s.C = s.A << s.Combo(operand)
	default:
		panic(fmt.Sprintf("bad instruction: %d", opcode))
	}

	s.pc += 2
}

func Parse(input io.Reader) State {
	state := State{}
	state.program = []int64{}
	state.output = []int64{}
	programstr := ""

	fmt.Fscanf(input, "Register A: %d\nRegister B: %d\nRegister C: %d\n\nProgram: %s", &state.A, &state.B, &state.C, &programstr)

	for _, instruction := range strings.Split(programstr, ",") {
		inst, _ := strconv.ParseInt(instruction, 10, 64)
		state.program = append(state.program, inst)
	}

	return state
}

func JoinCommas(nums []int64) string {
	str := []string{}

	for _, num := range nums {
		str = append(str, strconv.FormatInt(num, 10))
	}

	return strings.Join(str, ",")
}

func RunToHalt(state *State) {
	for {
		halt, opcode, operand := state.Fetch()

		if halt {
			break
		}

		state.Execute(opcode, operand)
	}
}

func Shiftall(nums []int64) int64 {
	a := int64(0)
	for _, num := range nums {
		a = (a << 3) + num
	}

	return a
}

func Unshiftall(hugenum int64) []int64 {
	nums := []int64{}

	for hugenum > 0 {
		nums = append(nums, hugenum%8)
		hugenum >>= 3
	}

	return nums
}

func main() {
	state := Parse(os.Stdin)
	a, b, c := state.A, state.B, state.C

	RunToHalt(&state)

	println(JoinCommas(state.output))

	result := Shiftall(state.output)
	resultfast := ExecFast(a, b, c)
	println(result)
	println(resultfast)
}
