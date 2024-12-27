package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	wire1, wire2 string
	gate         string
	outputwire   string
}

func (i Instruction) Ready(wires map[string]bool) bool {
	_, exists1 := wires[i.wire1]
	_, exists2 := wires[i.wire2]

	return exists1 && exists2
}

func (i *Instruction) Normalize() {
	if strings.Compare(i.wire1, i.wire2) > 0 {
		(*i).wire1, (*i).wire2 = i.wire2, i.wire1
	}
}

func (i Instruction) Apply(wires map[string]bool) {
	if !i.Ready(wires) {
		panic("not ready yet")
	}

	v1 := wires[i.wire1]
	v2 := wires[i.wire2]

	switch i.gate {
	case "AND":
		wires[i.outputwire] = v1 && v2
	case "OR":
		wires[i.outputwire] = v1 || v2
	case "XOR":
		wires[i.outputwire] = v1 != v2
	default:
		panic("bad instruction")
	}
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s %s %s -> %s", i.wire1, i.gate, i.wire2, i.outputwire)
}

func main() {
	inputfull, _ := io.ReadAll(os.Stdin)
	input := strings.Split(string(inputfull), "\n\n")

	wires := map[string]bool{}

	instructions := []Instruction{}

	for _, line := range strings.Split(strings.TrimSpace(input[0]), "\n") {
		linesegments := strings.Split(line, ": ")
		wire := linesegments[0]
		val, _ := strconv.ParseInt(linesegments[1], 10, 32)

		wires[wire] = val == 1
	}

	for _, line := range strings.Split(strings.TrimSpace(input[1]), "\n") {
		i := Instruction{}
		fmt.Sscanf(line, "%s %s %s -> %s", &i.wire1, &i.gate, &i.wire2, &i.outputwire)

		i.Normalize()

		instructions = append(instructions, i)
	}

	part1wires := maps.Clone(wires)
	ExecuteWhileChanged(part1wires, instructions)

	// DebugRenameCircuit(instructions)

	MakePlotCircuit(instructions)

	// DebugAdder(instructions)
	// PrintCircuit(instructions)

	println(MakeOutput(part1wires))
}

func ExecuteWhileChanged(wires map[string]bool, instructions []Instruction) {

	// create a map of changes to indicate which instructions we should re-run if the inputs have changed
	changes := map[string]struct{}{}

	for wire := range wires {
		changes[wire] = struct{}{}
	}

	// create a map to provide information on which change in wire can trigger each instruction
	wire2instructions := map[string][]Instruction{}

	for _, instruction := range instructions {
		if _, ok := wire2instructions[instruction.wire1]; !ok {
			wire2instructions[instruction.wire1] = []Instruction{}
		}

		if _, ok := wire2instructions[instruction.wire2]; !ok {
			wire2instructions[instruction.wire2] = []Instruction{}
		}

		wire2instructions[instruction.wire1] = append(wire2instructions[instruction.wire1], instruction)
		wire2instructions[instruction.wire2] = append(wire2instructions[instruction.wire2], instruction)
	}

	// execute all instructions that have changed while something changed in the loop
	for len(changes) > 0 {
		changed := map[string]struct{}{}

		for wire := range changes {
			ii := wire2instructions[wire]

			for _, inst := range ii {
				if !inst.Ready(wires) {
					continue
				}

				inst.Apply(wires)
				changed[inst.outputwire] = struct{}{}
			}
		}

		changes = changed
	}
}

func PerformAddition(instructions []Instruction, bitsize, x, y int) int64 {
	wires := map[string]bool{}

	maxint := 1<<(bitsize+1) - 1

	if x > maxint {
		panic("x overflow")
	}

	for i := bitsize; i >= 0; i-- {
		wirex := fmt.Sprintf("x%02d", i)
		wx := 1 << i
		wires[wirex] = (x&wx == wx)
	}

	if y > maxint {
		panic("y overflow")
	}

	for i := bitsize; i >= 0; i-- {
		wirey := fmt.Sprintf("y%02d", i)
		wy := 1 << i
		wires[wirey] = (x&wy == wy)
	}

	ExecuteWhileChanged(wires, instructions)

	// DebugRenameCircuit(instructions)

	return MakeOutput(wires)
}

func MakePlotCircuit(instructions []Instruction) {
	for _, inst := range instructions {
		fmt.Printf("[\"%s\"] -> [\"%s\"]\n", inst.wire1, inst.outputwire)
		fmt.Printf("[\"%s\"] -> [\"%s\"]\n", inst.wire2, inst.outputwire)
	}
}

func TestAddition(instructions []Instruction, bitsize, x, y int) {
	result := PerformAddition(instructions, bitsize-1, 1, 1)

	fmt.Printf("%d + %d = %d (%t)\n", x, y, result, (int64(x+y) == result))
}

func DebugAdder(instructions []Instruction) {
	// 1. determine the number of output gates
	numz := CountOutputGates(instructions)

	// for each output bit, generate a testcase that uniquely tests if the bits are adding up correctly
	// don't include the last bit (bit = numz) because then we overflow our adder
	for bit := 0; bit < numz; bit++ {
		TestAddition(instructions, numz, 1, 1<<bit)
		TestAddition(instructions, numz, 1<<bit, 1)
		TestAddition(instructions, numz, 1<<bit, 1<<bit)
	}
}

func FindInstruction(template Instruction, instructions []Instruction) (Instruction, bool) {
	for _, instruction := range instructions {
		if template.wire1 != "" && instruction.wire1 != template.wire1 {
			continue
		}

		if template.wire2 != "" && instruction.wire2 != template.wire2 {
			continue
		}

		if template.gate != "" && instruction.gate != template.gate {
			continue
		}

		if template.outputwire != "" && instruction.outputwire != template.outputwire {
			continue
		}

		return instruction, true
	}

	// try again with wire1 and wire2 swapped
	template.wire1, template.wire2 = template.wire2, template.wire1

	for _, instruction := range instructions {
		if template.wire1 != "" && instruction.wire1 != template.wire1 {
			continue
		}

		if template.wire2 != "" && instruction.wire2 != template.wire2 {
			continue
		}

		if template.gate != "" && instruction.gate != template.gate {
			continue
		}

		if template.outputwire != "" && instruction.outputwire != template.outputwire {
			continue
		}

		return instruction, true
	}

	return template, false
}

func RenameInstructions(instructions []Instruction, from, to string) []Instruction {
	result := []Instruction{}

	for _, instruction := range instructions {
		changed := instruction

		if changed.wire1 == from {
			changed.wire1 = to
		}

		if changed.wire2 == from {
			changed.wire2 = to
		}

		if changed.outputwire == from {
			changed.outputwire = to
		}

		changed.Normalize()

		result = append(result, changed)
	}

	return result
}

func MakeAdderCircuit(x, y, z string, instructions []Instruction) {
	// There should be a line for A xor B -> line1
	line1, found := FindInstruction(Instruction{wire1: y, wire2: x, gate: "XOR"}, instructions)

	if !found {
		fmt.Printf("%s XOR %s -> line1 not found\n", x, y)
		return
	}

	// There should be a line for carry xor line1 -> Z
	carry, found := FindInstruction(Instruction{wire2: line1.outputwire, gate: "XOR", outputwire: z}, instructions)

	if !found {
		fmt.Printf("carry XOR line1 (%s) -> %s not found\n", line1.outputwire, z)
		return
	}

	// There should be a line for line1 and carry -> line2
	line2, found := FindInstruction(Instruction{wire1: line1.outputwire, wire2: carry.outputwire, gate: "AND"}, instructions)

	if !found {
		fmt.Printf("line1 (%s) AND carry (%s) -> line2 not found\n", line1.outputwire, carry.outputwire)
		return
	}

	// There should be a line for A and B -> line3
	line3, found := FindInstruction(Instruction{wire1: x, wire2: y, gate: "AND"}, instructions)

	if !found {
		fmt.Printf("A AND B -> line3 not found\n")
		return
	}

	// There should be a line for line2 or line3 -> carry_out
	carryOut, found := FindInstruction(Instruction{wire1: line2.outputwire, wire2: line3.outputwire, gate: "OR"}, instructions)

	if !found {
		fmt.Printf("line2 OR line3 -> carry_out")
		return
	}

	println(line1.String())
	println(line2.String())
	println(line3.String())
	println(carry.String())
	println(carryOut.String())
}

func DebugRenameCircuit(instructions []Instruction) {
	instructions2 := instructions
	// rename the output of x00 AND y00 to 00AND00
	for _, inst := range instructions {
		if inst.wire1[0] == 'x' && inst.wire2[0] == 'y' {
			instructions2 = RenameInstructions(instructions2, inst.outputwire, fmt.Sprintf("%s%s%s", inst.wire1, inst.gate, inst.wire2))
		}
	}

	instructions = instructions2

	for _, inst := range instructions {
		fmt.Println(inst.String())
	}
}

func PrintCircuit(instructions []Instruction) {
	numz := CountOutputGates(instructions)

	for i := 0; i < numz; i++ {
		println(i)
		x := fmt.Sprintf("%c%02d", 'x', i)
		y := fmt.Sprintf("%c%02d", 'y', i)
		z := fmt.Sprintf("%c%02d", 'z', i)
		MakeAdderCircuit(x, y, z, instructions)
	}
}

func InitializeWires(wires *map[string]bool, totalbits int, name rune, bit int, bitvalue bool) int64 {
	for i := 0; i <= totalbits; i++ {
		wirename := fmt.Sprintf("%c%2d", name, i)
		(*wires)[wirename] = bitvalue && i == bit
	}

	if bitvalue == false {
		return 0
	}

	return 1 << bit
}

func CountOutputGates(instructions []Instruction) int {
	numoutputs := int64(0)
	for _, instruction := range instructions {
		if instruction.outputwire[0] == 'z' {
			id, _ := strconv.ParseInt(instruction.outputwire[1:], 10, 32)
			fmt.Printf("found zID: %s (%d)\n", instruction.outputwire, id)
			numoutputs = max(numoutputs, id)
		}
	}

	return int(numoutputs)
}

func MakeOutput(wires map[string]bool) int64 {
	num := int64(0)

	for wire, value := range wires {
		if !value {
			continue
		}

		if wire[0] == 'z' {
			bitshift, _ := strconv.ParseInt(wire[1:], 10, 32)

			num |= 1 << bitshift
		}
	}

	return num
}
