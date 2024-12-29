package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
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

func (i Instruction) HasWire(wire string) bool {
	return i.wire1 == wire || i.wire2 == wire
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
	// MakePlotCircuit(instructions)
	// DebugAdder(instructions)
	// PrintCircuit(instructions)
	// FindAllPathsToZ(instructions)
	// StructuredSwapping(instructions)
	// DetectBadGraph(instructions)
	Swap4(instructions)

	rootinstructions := []Instruction{}

	for _, inst := range instructions {
		if inst.wire1[0] == 'x' || inst.wire1[0] == 'Y' {
			rootinstructions = append(rootinstructions, inst)
		}
	}

	// DetectAnomalies(instructions, rootinstructions)

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

func PerformAddition(instructions []Instruction, bitsize, x, y int64) int64 {
	wires := map[string]bool{}

	maxint := int64(1<<(bitsize+1) - 1)

	if x > maxint {
		panic("x overflow")
	}

	for i := bitsize; i >= 0; i-- {
		wirex := fmt.Sprintf("x%02d", i)
		wx := int64(1 << i)
		wires[wirex] = (x&wx == wx)
	}

	if y > maxint {
		panic("y overflow")
	}

	for i := bitsize; i >= 0; i-- {
		wirey := fmt.Sprintf("y%02d", i)
		wy := int64(1 << i)
		wires[wirey] = (x&wy == wy)
	}

	ExecuteWhileChanged(wires, instructions)

	return MakeOutput(wires)
}

func Bfs(instructions []Instruction, start Instruction, endwire string) []Instruction {
	wire2instruction := map[string][]Instruction{}

	for _, inst := range instructions {
		if wire2instruction[inst.wire1] == nil {
			wire2instruction[inst.wire1] = []Instruction{}
		}
		if wire2instruction[inst.wire2] == nil {
			wire2instruction[inst.wire2] = []Instruction{}
		}

		wire2instruction[inst.wire1] = append(wire2instruction[inst.wire1], inst)
		wire2instruction[inst.wire2] = append(wire2instruction[inst.wire2], inst)
	}

	prev := map[Instruction]Instruction{}
	queue := []Instruction{start}

	for len(queue) > 0 {
		element := queue[0]

		if element.outputwire == endwire {
			break
		}

		queue = queue[1:]

		for _, neighbour := range wire2instruction[element.outputwire] {
			prev[neighbour] = element
			queue = append(queue, neighbour)
		}
	}

	path := []Instruction{}

	if len(queue) == 0 {
		// didnt exit the loop via the break statement, no path found
		return path
	}
	p := queue[0]

	for p != (Instruction{}) {
		path = append(path, p)
		p = prev[p]
	}

	slices.Reverse(path)

	return path
}

type GatePattern struct {
	xor, and, or int
}

func DetectAnomalies(instructions, current []Instruction) {
	next := map[Instruction]struct{}{}
	wire2instruction := map[string][]Instruction{}

	// patterns[x] = y means instruction x feeds into pattern y, e.g. 2 XORs and 1 OR
	patterns := map[Instruction]GatePattern{}

	for _, inst := range instructions {
		if wire2instruction[inst.wire1] == nil {
			wire2instruction[inst.wire1] = []Instruction{}
		}
		if wire2instruction[inst.wire2] == nil {
			wire2instruction[inst.wire2] = []Instruction{}
		}

		wire2instruction[inst.wire1] = append(wire2instruction[inst.wire1], inst)
		wire2instruction[inst.wire2] = append(wire2instruction[inst.wire2], inst)
	}

	for _, inst := range current {
		// detect what sort of destinations this instruction has
		for _, destination := range wire2instruction[inst.outputwire] {
			pattern := patterns[destination]
			if destination.gate == "XOR" {
				pattern.xor++
			} else if destination.gate == "AND" {
				pattern.and++
			} else if destination.gate == "OR" {
				pattern.or++
			}

			patterns[destination] = pattern
			next[destination] = struct{}{}
		}
	}

	// expect that at least half of the gates are correct, so if a pattern
	// occured fewer times, there is probably a candiate for corruption.
	counts, inverted := CountAnomalies(patterns)
	for pattern, count := range counts {
		if count > len(current)/4 {
			fmt.Printf("regular: %+v (%d times)\n", pattern, count)
			continue
		}

		fmt.Printf("found anomaly pattern %+v %d times on gates\n", pattern, count)

		for _, instruction := range inverted[pattern] {
			fmt.Printf("\t%s\n", instruction.String())
		}
	}

	if len(next) > 0 {
		nextitems := slices.Collect(maps.Keys(next))
		DetectAnomalies(instructions, nextitems)
	}
}

func CountAnomalies[K, V comparable](cases map[K]V) (map[V]int, map[V][]K) {
	result := map[V][]K{}
	counts := map[V]int{}

	for key, value := range cases {
		counts[value]++
		if _, exists := result[value]; !exists {
			result[value] = []K{}
		}

		result[value] = append(result[value], key)
	}

	return counts, result
}

func FindAllPathsToZ(instructions []Instruction) int {
	slices.SortFunc(instructions, func(a, b Instruction) int {
		return strings.Compare(a.String(), b.String())
	})

	outputs := int(CountOutputGates(instructions))

	for bit := 0; bit < outputs; bit++ {
		x, y, z, z2 := fmt.Sprintf("x%02d", bit), fmt.Sprintf("y%02d", bit), fmt.Sprintf("z%02d", bit), fmt.Sprintf("z%02d", bit+1)

		// x XOR y -> z
		xor, _ := FindInstruction(Instruction{wire1: x, wire2: y, gate: "XOR"}, instructions)
		and, _ := FindInstruction(Instruction{wire1: x, wire2: y, gate: "AND"}, instructions)

		path1 := Bfs(instructions, xor, z)

		// fmt.Printf("%s XOR %s", x, y)

		// for _, p := range path1 {
		// 	fmt.Printf(" -> %s", p.gate)
		// }

		// fmt.Printf("\n%s AND %s", x, y)

		path2 := Bfs(instructions, and, z2)

		// for _, p := range path2 {
		// 	fmt.Printf(" -> %s", p.gate)
		// }
		//
		// fmt.Printf("\n")

		// from wikipedia:
		// Sum = A xor B xor C_in
		// C_out = (A and B) or (C_in and (A xor B))
		//
		// let A xor B = halfsum
		// let A and B = firstcarry
		// let C_in and halfsum = secondcarry

		// ->
		// halfsum     = x XOR y
		// z_n         = halfsum XOR c_in
		// firstcarry  = x AND y
		// secondcarry = halfsum AND C_in
		// C_out       = firstcarry OR secondcarry

		expected1 := []string{"XOR", "XOR"}
		expected2 := []string{"AND", "OR", "XOR"}

		if bit == 0 {
			expected1 = []string{"XOR"}
			expected2 = []string{"AND", "XOR"}
		}

		gates1, gates2 := []string{}, []string{}

		for _, p := range path1 {
			gates1 = append(gates1, p.gate)
		}

		for _, p := range path2 {
			gates2 = append(gates2, p.gate)
		}

		if !slices.Equal(gates2, expected1) || !slices.Equal(gates2, expected2) {
			return bit
		}
	}

	return outputs
}

type Swap struct {
	from, to *string
}

func StructuredSwapping(instructions []Instruction) {
	slices.SortFunc(instructions, func(a, b Instruction) int {
		return strings.Compare(a.String(), b.String())
	})

	outputs := []*string{}

	for _, inst := range instructions {
		outputs = append(outputs, &inst.outputwire)
	}

	// find pair that fixes most broken bits

	// first pair
	firstbits, firstswap, secondswap := Try2Swaps(instructions, outputs)
	fmt.Printf("first: %d bits correct\n", firstbits)
	*firstswap.to, *firstswap.from = *firstswap.from, *firstswap.to
	*secondswap.to, *secondswap.from = *secondswap.from, *secondswap.to

	// third pair
	thirdbits, thirdswap := TrySwaps(instructions, outputs)
	fmt.Printf("third: %d bits correct\n", thirdbits)
	*thirdswap.to, *thirdswap.from = *thirdswap.from, *thirdswap.to

	// fourth pair
	fourthbits, fourthswap := TrySwaps(instructions, outputs)
	fmt.Printf("fourth: %d bits correct\n", fourthbits)
	*fourthswap.to, *fourthswap.from = *fourthswap.from, *fourthswap.to
}

func TrySwaps(instructions []Instruction, outputs []*string) (int, Swap) {
	correctbits := 0
	swap := Swap{}

	for i, from := range outputs {
		for _, to := range outputs[i+1:] {
			*from, *to = *to, *from
			_, bits := FindFirstBrokenBit(instructions)
			if bits > correctbits {
				correctbits = bits
				swap = Swap{from, to}
			}
			*from, *to = *to, *from
		}
	}

	return correctbits, swap
}

func Try2Swaps(instructions []Instruction, outputs []*string) (int, Swap, Swap) {
	correctbits := 0
	swap1, swap2 := Swap{}, Swap{}
	for i, from := range outputs {
		for j, to := range outputs[i+1:] {
			for k, from2 := range outputs[i+j+1:] {
				for _, to2 := range outputs[i+j+k+1:] {

					*from, *to = *to, *from
					*from2, *to2 = *to2, *from2
					_, bits := FindFirstBrokenBit(instructions)
					if bits > correctbits {
						correctbits = bits
						swap1 = Swap{from, to}
						swap2 = Swap{from2, to2}
					}

					*from, *to = *to, *from
					*from2, *to2 = *to2, *from2
				}
			}
		}
	}

	return correctbits, swap1, swap2
}

func CountMistakes(instructions []Instruction, maxbits int) int {
	// test 001000 + 0 where 1 moves over 1 bit
	mistakes := 0
	bitsize := min(CountOutputGates(instructions), int64(maxbits))
	for offset := int64(0); offset < bitsize; offset++ {
		largenum := int64(1 << offset)
		addition := PerformAddition(instructions, bitsize, largenum, 0)

		if largenum != addition {
			mistakes++
			// fmt.Printf("%0*s + 0 != %0*s (offset = %d)\n", bitsize, strconv.FormatInt(largenum, 2), bitsize, strconv.FormatInt(addition, 2), offset)
		}

		addition = PerformAddition(instructions, bitsize, 0, largenum)

		if largenum != addition {
			mistakes++
			// fmt.Printf("0 + %0*s != %0*s (offset = %d)\n", bitsize, strconv.FormatInt(largenum, 2), bitsize, strconv.FormatInt(addition, 2), offset)
		}
	}

	// test 111111 + 0, adding 1 extra bit each time
	for offset := int64(0); offset < bitsize; offset++ {
		largenum := int64(1<<(offset+1) - 1)
		addition := PerformAddition(instructions, bitsize, largenum, 0)

		if largenum != addition {
			mistakes++
			// fmt.Printf("%0*s + 0 != %0*s (offset = %d)\n", bitsize, strconv.FormatInt(largenum, 2), bitsize, strconv.FormatInt(addition, 2), offset)
		}

		addition = PerformAddition(instructions, bitsize, 0, largenum)

		if largenum != addition {
			mistakes++
			// fmt.Printf("0 + %0*s != %0*s (offset = %d)\n", bitsize, strconv.FormatInt(largenum, 2), bitsize, strconv.FormatInt(addition, 2), offset)
		}
	}

	return mistakes
}

func MakePlotCircuit(instructions []Instruction) {
	for _, inst := range instructions {
		fmt.Printf("[\"%s\"] -> [\"%s\"]\n", inst.wire1, inst.outputwire)
		fmt.Printf("[\"%s\"] -> [\"%s\"]\n", inst.wire2, inst.outputwire)
	}
}

func TestAddition(instructions []Instruction, bitsize int64, x, y int) {
	result := PerformAddition(instructions, bitsize-1, 1, 1)

	fmt.Printf("%d + %d = %d (%t)\n", x, y, result, (int64(x+y) == result))
}

func FindFirstBrokenBit(instructions []Instruction) (bool, int) {
	bitsize := CountOutputGates(instructions)

	x := int64(1<<(bitsize) - 1)
	y := int64(1)
	zexpected := x + y
	zactual := PerformAddition(instructions, bitsize, x, y)

	for bit := 0; bit < int(bitsize); bit++ {
		if zexpected&(1<<bit) != zactual&(1<<bit) {
			return false, bit
		}
	}

	return true, 0
}

func DebugAdder(instructions []Instruction) {
	// 1. determine the number of output gates
	numz := CountOutputGates(instructions)

	// for each output bit, generate a testcase that uniquely tests if the bits are adding up correctly
	// don't include the last bit (bit = numz) because then we overflow our adder
	for bit := int64(0); bit < numz; bit++ {
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

func FindInstructions(template Instruction, instructions []Instruction) []Instruction {
	result := []Instruction{}

	// ensure wire order is the same
	template.Normalize()
	for _, instruction := range instructions {
		instruction.Normalize()

		if template.wire1 != "" && instruction.wire1 != template.wire1 && instruction.wire2 != template.wire1 {
			continue
		}

		if template.wire2 != "" && instruction.wire2 != template.wire2 && instruction.wire1 != template.wire2 {
			continue
		}

		if template.gate != "" && instruction.gate != template.gate {
			continue
		}

		if template.outputwire != "" && instruction.outputwire != template.outputwire {
			continue
		}

		result = append(result, instruction)
	}

	return result
}

func SwapOutputs(instructions *[]Instruction, from, to string) {
	for i := range *instructions {
		wo := (*instructions)[i].outputwire

		if wo == from {
			(*instructions)[i].outputwire = to
		} else if wo == to {
			(*instructions)[i].outputwire = from
		}
	}
}

func RenameInstructions(instructions *[]Instruction, from, to string) {
	for i, inst := range *instructions {
		if inst.wire1 == from {
			inst.wire1 = to
		}

		if inst.wire2 == from {
			inst.wire2 = to
		}

		if inst.outputwire == from {
			inst.outputwire = to
		}

		inst.Normalize()

		(*instructions)[i] = inst
	}
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
	substitutions := map[string]string{}

	// from wikipedia:
	// Sum = A xor B xor C_in
	// C_out = (A and B) or (C_in and (A xor B))
	//
	// let A xor B = halfsum
	// let A and B = firstcarry
	// let C_in and halfsum = secondcarry

	// ->
	// halfsum     = x XOR y
	// z_n         = halfsum XOR c_in
	// firstcarry  = x AND y
	// secondcarry = halfsum AND C_in
	// C_out       = firstcarry OR secondcarry

	// rename the output of x00 AND y00 to 00AND00
	for _, inst := range instructions {
		if inst.wire1[0] == 'x' && inst.wire2[0] == 'y' {
			number, _ := strconv.ParseInt(inst.wire1[1:], 10, 32)

			if inst.gate == "XOR" {
				// x01 XOR y01 should be halfsum01
				// x01 XOR y01 XOR c00 should be fullsum
				to := fmt.Sprintf("halfsum%02d", number)
				substitutions[to] = inst.outputwire
				RenameInstructions(&instructions, inst.outputwire, to)
			}

			if inst.gate == "AND" {
				// x01 AND y01 is the half carry for 02, which should be combined with C_in && (A xor B)
				to := fmt.Sprintf("firstcarry%02d", number)
				substitutions[to] = inst.outputwire
				RenameInstructions(&instructions, inst.outputwire, to)
			}
		}
	}

	for _, inst := range instructions {
		if inst.outputwire[0] == 'z' && inst.gate == "XOR" {
			// this is the instruction that adds the previous carry to the halfsum
			// the name of the halfsum has already been renamed in the loop above
			nothalfsum := inst.wire1

			if len(inst.wire1) > len(inst.wire2) {
				nothalfsum = inst.wire2
			}

			number, _ := strconv.ParseInt(inst.outputwire[1:], 10, 32)
			to := fmt.Sprintf("carry_in%02d", number)

			substitutions[to] = nothalfsum

			RenameInstructions(&instructions, nothalfsum, to)
		}

		longest, shortest := inst.wire1, inst.wire2

		if len(longest) < len(shortest) {
			shortest, longest = longest, shortest
		}

		if inst.gate == "OR" && strings.HasPrefix(longest, "firstcarry") {
			var num int
			fmt.Sscanf(longest, "firstcarry%02d", &num)

			RenameInstructions(&instructions, shortest, fmt.Sprintf("secondcarry%02d", num))
		}
	}

	for _, inst := range instructions {
		inst.Normalize()
		fmt.Println(inst.String())
	}

	interested := [...]string{"firstcarry00", "carry_in01", "firstcarry11", "halfsum01", "pvh"}
	replacements := []string{}

	for _, s := range interested {
		fmt.Printf("zz %s -> %s\n", s, substitutions[s])
		if sub, exists := substitutions[s]; exists {
			replacements = append(replacements, sub)
		} else {
			replacements = append(replacements, s)
		}
	}

	// NaiveSwapping(instructions, []int{}, replacements)
}

func Swap4(instructions []Instruction) {
	baseline := VerifyAll(instructions)

	swaps := []string{}

	outputs := []*string{}

	for i := range instructions {
		outputs = append(outputs, &instructions[i].outputwire)
	}

	fmt.Printf("baseline: %d\n", baseline)

	for swap := 0; swap < 4; swap++ {
		// swap every element, and check for improvement
	swapall:
		for i, from := range outputs {
			for _, to := range outputs[i+1:] {
				*from, *to = *to, *from
				correct := VerifyAll(instructions)

				if correct > baseline {
					baseline = correct
					swaps = append(swaps, *from, *to)
					break swapall
				}

				*from, *to = *to, *from
			}
		}
	}

	fmt.Printf("%+v\n", swaps)
	slices.Sort(swaps)
	swaps = slices.Compact(swaps)

	fmt.Printf("%s\n", strings.Join(swaps, ","))
}

// from wikipedia:
// Sum = A xor B xor C_in
// C_out = (A and B) or (C_in and (A xor B))
//
// let A xor B = halfsum
// let A and B = firstcarry
// let C_in and halfsum = secondcarry

// ->
// halfsum     = x XOR y
// z_n         = halfsum XOR c_in
// firstcarry  = x AND y
// secondcarry = halfsum AND C_in
// C_out       = firstcarry OR secondcarry
func VerifyOutput(instructions map[string]Instruction, bit int) bool {
	// the output comes from an XOR. If it's the first bit, it should be an XOR
	// of the inputs, if not it should be an XOR of the XOR of the first bits,
	// and the carry
	z := fmt.Sprintf("z%02d", bit)
	x, y := fmt.Sprintf("x%02d", bit), fmt.Sprintf("y%02d", bit)
	parent := instructions[z]

	// fmt.Printf("start verifying output %s\n", parent.String())

	if parent.gate != "XOR" {
		return false
	}

	if bit == 0 {
		return parent.HasWire(x) && parent.HasWire(y)
	}

	leftislowercarry := VerifyLowerCarry(instructions, parent.wire1, bit-1)
	rightislowercarry := VerifyLowerCarry(instructions, parent.wire2, bit-1)
	leftishalfsum := VerifyHalfsum(instructions, parent.wire1, bit)
	rightishalfsum := VerifyHalfsum(instructions, parent.wire2, bit)

	// fmt.Printf("%t && %t or %t && %t\n", leftislowercarry, rightishalfsum, leftishalfsum, rightislowercarry)

	return leftislowercarry && rightishalfsum || leftishalfsum && rightislowercarry

	// here we validate that one of the wires is a halfsum, and the other half is the lower carry
	if !(VerifyHalfsum(instructions, parent.wire1, bit) && VerifyLowerCarry(instructions, parent.wire2, bit-1)) || !(VerifyHalfsum(instructions, parent.wire2, bit) && VerifyLowerCarry(instructions, parent.wire1, bit-1)) {

		fmt.Printf("verifying %s is wrong\n", parent.String())
		return false
	}

	return true
}

func VerifyHalfsum(instructions map[string]Instruction, wire string, bit int) bool {
	// this is the product of x XOR y
	parent := instructions[wire]

	// fmt.Printf("expecting halfsum %s\n", parent.String())

	if parent.gate != "XOR" {
		// fmt.Printf("expected XOR, got %s\n", parent.gate)
		return false
	}

	x, y := fmt.Sprintf("x%02d", bit), fmt.Sprintf("y%02d", bit)

	return parent.HasWire(x) && parent.HasWire(y)
}

func VerifyLowerCarry(instructions map[string]Instruction, wire string, bit int) bool {
	// on bit != 0 this will be the product of (A and B) OR (C_in and (A xor B))
	// on bit == 0 this will just be (A and B)

	parent := instructions[wire]

	if bit == 0 {
		// fmt.Printf("verifying lower carry bit: %d\n", bit)
		// should be A and B
		return VerifyFirstCarry(instructions, wire, bit)
	}

	// fmt.Printf("expecting %s to be of shape (firstcarry) OR (secondcarry)\n", parent.String())

	// this should be an OR of firstcarry OR secondcarry
	if parent.gate != "OR" {
		return false
	}

	return VerifyFirstCarry(instructions, parent.wire1, bit) && VerifySecondCarry(instructions, parent.wire2, bit) || VerifyFirstCarry(instructions, parent.wire2, bit) && VerifySecondCarry(instructions, parent.wire1, bit)
}

func VerifyFirstCarry(instructions map[string]Instruction, wire string, bit int) bool {
	// x AND y
	parent := instructions[wire]
	// fmt.Printf("expecting %s to be of shape x AND y: %d\n", parent.String(), bit)
	x, y := fmt.Sprintf("x%02d", bit), fmt.Sprintf("y%02d", bit)

	if parent.gate != "AND" {
		// fmt.Printf("expected AND, got %s\n", parent.gate)
		return false
	}

	return parent.HasWire(x) && parent.HasWire(y)
}

func VerifySecondCarry(instructions map[string]Instruction, wire string, bit int) bool {
	// halfsum AND c_in
	parent := instructions[wire]

	if parent.gate != "AND" {
		return false
	}

	return VerifyHalfsum(instructions, parent.wire1, bit) && VerifyLowerCarry(instructions, parent.wire2, bit-1) || VerifyHalfsum(instructions, parent.wire2, bit) && VerifyLowerCarry(instructions, parent.wire1, bit-1)
}

func VerifyAll(instructions []Instruction) int {
	output2gates := map[string]Instruction{}

	for _, inst := range instructions {
		output2gates[inst.outputwire] = inst
	}

	total := int(CountOutputGates(instructions))

	for bit := 0; bit < total; bit++ {
		if !VerifyOutput(output2gates, bit) {
			return bit
		}
	}

	return total
}

func DetectBadGraph(instructions []Instruction) []string {
	// from wikipedia:
	// Sum = A xor B xor C_in
	// C_out = (A and B) or (C_in and (A xor B))
	//
	// let A xor B = halfsum
	// let A and B = firstcarry
	// let C_in and halfsum = secondcarry

	// ->
	// halfsum     = x XOR y
	// z_n         = halfsum XOR c_in
	// firstcarry  = x AND y
	// secondcarry = halfsum AND C_in
	// C_out       = firstcarry OR secondcarry

	// for xN, yN: N: 0.., this should go into an XOR, and into an AND instruction. if it doesn't the output is wrong.

	bad := []string{}

	for n := 0; n < int(CountOutputGates(instructions)); n++ {
		xname, yname, zname := fmt.Sprintf("x%02d", n), fmt.Sprintf("y%02d", n), fmt.Sprintf("z%02d", n)

		halfsum, found := FindInstruction(Instruction{wire1: xname, wire2: yname, gate: "XOR"}, instructions)

		if !found {
			panic("x XOR y not found")
		}

		firstcarry, found := FindInstruction(Instruction{wire1: xname, wire2: yname, gate: "AND"}, instructions)

		if !found {
			panic("x AND y not found")
		}

		fullcarry1 := FindInstructions(Instruction{wire1: firstcarry.outputwire, gate: "OR"}, instructions)

		if len(fullcarry1) != 1 {
			fmt.Printf("found bad output: %s (should go to 1 OR output, goes to %d)\n", firstcarry.outputwire, len(fullcarry1))
			bad = append(bad, firstcarry.outputwire)
		}

		// for the unique case of n == 0, the XOR goes into the z instruction directly.
		if n == 0 {
			if halfsum.outputwire[0] != 'z' {
				fmt.Printf("found bad output: %s\n", halfsum.outputwire)
				bad = append(bad, halfsum.outputwire)
				continue
			}

			next := FindInstructions(Instruction{wire1: halfsum.outputwire, gate: "AND"}, instructions)

			if len(next) != 1 {
				fmt.Printf("found bad output: %s\n", halfsum.outputwire)
				bad = append(bad, halfsum.outputwire)
				continue
			}
		} else {
			// this output should go to 2 instructions: halfsum XOR c_in, halfsum AND c_in
			next := FindInstructions(Instruction{wire1: halfsum.outputwire, gate: "XOR"}, instructions)

			if len(next) != 1 {
				fmt.Printf("found bad output: %s (%d other outputs)\n", halfsum.outputwire, len(next))
				bad = append(bad, halfsum.outputwire)
				continue
			}

			// the output of this one should be the bit of the number
			if next[0].outputwire != zname {
				fmt.Printf("found bad output: %s (should go to %s)\n", next[0].outputwire, zname)
				bad = append(bad, next[0].outputwire)
				continue
			}

			c_in := next[0].wire2

			if c_in == halfsum.outputwire {
				c_in = next[0].wire1
			}

			next = FindInstructions(Instruction{wire1: halfsum.outputwire, wire2: c_in, gate: "AND"}, instructions)

			if len(next) != 1 {
				fmt.Printf("found bad output: %s (should go to 1 output, goes to %d outputs)\n", halfsum.outputwire, len(next))
				bad = append(bad, halfsum.outputwire)
				continue
			}

			fullcarry := FindInstructions(Instruction{wire1: next[0].outputwire, gate: "OR"}, instructions)

			if len(fullcarry) != 1 {
				fmt.Printf("found bad output: %s, should go to OR, goes to %d gates", next[0].outputwire, len(fullcarry))
				bad = append(bad, next[0].outputwire)
				continue
			}
		}
	}

	return bad
}

func NaiveSwapping(instructions []Instruction, used []int, available []string) {
	if len(used) == 8 {
		usedstrings := []string{}
		for i := 0; i < 8; i += 2 {
			usedstrings = append(usedstrings, available[used[i]], available[used[i+1]])
			SwapOutputs(&instructions, available[used[i]], available[used[i+1]])
		}
		mistakes := CountMistakes(instructions, 99)
		fmt.Printf("pairs: %+v -> mistakes: %d\n", used, mistakes)

		if mistakes == 0 {
			slices.Sort(used)
			fmt.Printf("%s\n", strings.Join(usedstrings, ","))
		}
		for i := 0; i < 8; i += 2 {
			SwapOutputs(&instructions, available[used[i]], available[used[i+1]])
		}

		return
	}

pickanumber:
	for i := range available {
		for _, u := range used {
			if i == u {
				continue pickanumber
			}
		}

	pickanothernumber:
		for j := i + 1; j < len(available); j++ {
			for _, u := range used {
				if i == u {
					continue pickanothernumber
				}
			}
			newused := []int{}
			newused = append(newused, used...)
			newused = append(newused, i, j)
			NaiveSwapping(instructions, newused, available)
		}
	}
}

func PrintCircuit(instructions []Instruction) {
	numz := CountOutputGates(instructions)

	for i := int64(0); i < numz; i++ {
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

func CountOutputGates(instructions []Instruction) int64 {
	numoutputs := int64(0)
	for _, instruction := range instructions {
		if instruction.outputwire[0] == 'z' {
			id, _ := strconv.ParseInt(instruction.outputwire[1:], 10, 32)
			numoutputs = max(numoutputs, id)
		}
	}

	return numoutputs
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
