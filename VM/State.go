package VM

type dataType uint8

const (
	String  dataType = 0
	Integer dataType = 1
)

//type Frame struct {
//	pc             uint
//	localVariables []*OracleListener.BroadcastMsg
//}

type State struct {
	Stack       *Stack
	Memory      Memory
	pc          uint
	consumedGas uint64
}

func newState() *State {

	return &State{
		Stack:       newStack(),
		Memory:      newMemory(),
		pc:          0,
		consumedGas: 0,
	}
}
