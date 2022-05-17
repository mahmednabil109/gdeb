package VM

import "strconv"

type dataType uint8

const (
	String  dataType = 0
	Integer dataType = 1
)

type Message struct {
	val      string
	dataType dataType
}
type Frame struct {
	Stack          *Stack
	pc             uint
	localVariables []chan Message
}

func newFrame() *Frame {
	return &Frame{
		Stack:          newStack(),
		pc:             0,
		localVariables: []chan Message{},
	}
}

type VMState struct {
	Memory      Memory
	Frame       *Frame
	consumedGas uint64
}

func newVM() *VMState {

	return &VMState{
		Memory:      newMemory(),
		Frame:       newFrame(),
		consumedGas: 0,
	}
}

func (vm *VMState) toString() string {
	return vm.Frame.Stack.toString() +
		"\n" + "PC ---->" + strconv.Itoa(int(vm.Frame.pc)) +
		"\n" + "Consumed Gas ----->" + strconv.Itoa(int(vm.consumedGas))
}
