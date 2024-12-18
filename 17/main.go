package main

import (
	"fmt"
	"io"
	"os"
	"slices"

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

func ExecuteProgram(program []int64, a int64) []int64 {
	state := State{
		A:       a,
		program: program,
		output:  []int64{},
	}

	RunToHalt(&state)

	return state.output
}

// This is a beautiful piece of legacy code from when Mansur and I were at the
// office until 22:00 to try and figure this out. We manually compiled the
// program into Go, and managed to reduce it into this expression, which we
// hoped was invertible by hand. It turns out that ChatGPT o1 did manage to
// create a function that could invert this using backtracking, and it does
// indeed produce a correct result. However, I have no idea what it does.
//
// Part of my ideas was to brute force it, which is also why this simpler
// expression is nice. Instead of an output array it produces a single i64
// number which can be converted to- and from the expected output array using
// the Shiftall and Unshiftall functions.
func ExecFast(a, b, c int64) int64 {
	out := int64(0)

	// this is my input program, translated into Go
	for a > 0 {
		// b = a & 0b111 ^ 2
		// c = a >> b
		// b = b ^ c
		// b = b ^ 3
		out = out << 3
		out = (a & 0b111) ^ ((a>>((a&0b111)^2))^1)&0b111
		a = a >> 3
	}

	return out
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

// This is a search function I've made after browsing some other solutions, and
// coming to the realisation that one of my previous assumptions was reversed.
// During the engineering session with Mansur I erroneously concluded that
// adding extra bits at the least-signficant end of A could potentially
// influence future outputs due to leftover state being present in the device
// memory, and that for this reason it was not possible to construct a solution
// on the basis of optimal substructure. This is not entirely correct, because
// the solution (the A number) is generated in reverse. This means that bits
// that are prepended on A will affect future solutions instead of past
// solutions. That's not a problem because those future solutions are unchecked
// anyway, and we can add any other bitcombination to make it work somehow.
//
// That's exactly what this searcher method does: it recursively adds 3 bits
// onto the a number, as each triplet corresponds to an output digit. This way
// it progressively constructs the output, starting at the tail and moving to
// the head. The first element of the output comes from the least significant 3
// bits of A.
func FindValidAs(as, program []int64) []int64 {
	next := []int64{}

	for _, a := range as {
		for i := 0; i < 8; i++ {
			na := (a << 3) | int64(i)
			output := ExecuteProgram(program, na)

			if len(output) > len(program) {
				// in this case, we are producing programs that are too large,
				// which means in the last iteration of this function we've
				// reached our goal and we can return the a's we found thus
				// far.
				return as
			}

			if !slices.Equal(output, program[len(program)-len(output):]) {
				continue
			}

			next = append(next, na)
		}
	}

	if len(next) == 0 {
		panic("no candidates found")
	}

	return FindValidAs(next, program)
}

func main() {
	state := Parse(os.Stdin)
	RunToHalt(&state)

	println(JoinCommas(state.output))
	println(slices.Min(FindValidAs([]int64{0}, state.program)))
}
