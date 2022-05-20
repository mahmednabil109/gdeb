package VM

import (
	"github.com/mahmednabil109/gdeb/OracleConnection"
	"strconv"
)

type dataType uint8

const (
	String  dataType = 0
	Integer dataType = 1
)

type Frame struct {
	Stack          *Stack
	pc             uint
	localVariables []OracleConnection.BroadcastMsg
	buffer         []*OracleConnection.SubscribeMsg
}

func newFrame() *Frame {
	return &Frame{
		Stack: newStack(),
		pc:    0,
	}
}

type State struct {
	Memory      Memory
	Frame       *Frame
	consumedGas uint64
	OracleConnection.OraclePool
}

func newVM() *State {

	return &State{
		Memory:      newMemory(),
		Frame:       newFrame(),
		consumedGas: 0,
	}
}

func (vm *State) toString() string {
	return vm.Frame.Stack.toString() +
		"\n" + "PC ---->" + strconv.Itoa(int(vm.Frame.pc)) +
		"\n" + "Consumed Gas ----->" + strconv.Itoa(int(vm.consumedGas))
}
