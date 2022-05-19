package VM

import (
	"github.com/mahmednabil109/gdeb/OracleConnection"
)

type (
	OperationType func(*State, *GlobalData) error
)

func AddOP(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.Add(b))

	return nil
}

func SubOP(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.Sub(b))

	return nil
}

func MulOP(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	result, _ := a.Multiply(b)
	stack.Push(result)

	return nil
}

//GreaterOp Return 1 if a > b
func GreaterOp(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	isGreater := a.GT(b)

	var result DataWord

	if isGreater {
		result.SetUint32(1, 0)
	} else {
		result.SetUint32(0, 0)
	}

	stack.Push(result)

	return nil
}

func XorOP(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.Xor(b))

	return nil
}
func AndOP(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.And(b))

	return nil
}
func OrOP(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.Or(b))
	return nil
}

func NotOP(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack

	a := stack.Pop()

	stack.Push(a.Not())
	return nil
}

func PushOp(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack
	newData := NewDataWord()
	newData.SetDataWord(globalData.ContractCode[state.Frame.pc+1 : state.Frame.pc+33])
	stack.Push(newData)
	return nil

}

func PopOp(state *State, globalData *GlobalData) error {
	state.Frame.Stack.Pop()
	return nil

}

func MStoreOp(state *State, globalData *GlobalData) error {
	mem := state.Memory

	offset, value := state.Frame.Stack.Pop().toInt32(), state.Frame.Stack.Pop().toByteArray()
	mem.store(int(offset), value)
	return nil

}

func MLoadOp(state *State, globalData *GlobalData) error {
	mem := state.Memory
	offset := state.Frame.Stack.Pop().toInt32()
	size := state.Frame.Stack.Pop().toInt32()

	data, err := mem.load(int(offset), int(size))

	if err != nil {
		return err
	}
	x := NewDataWord()
	x.SetDataWord(data)
	state.Frame.Stack.Push(x)
	return nil
}

func JumpOp(state *State, globalData *GlobalData) error {
	pc := &state.Frame.pc
	*pc = uint(state.Frame.Stack.Pop().toInt32())
	return nil
}

// JumpIOp conditional Jump
func JumpIOp(state *State, globalData *GlobalData) error {
	pc := &state.Frame.pc
	check := state.Frame.Stack.Pop().toInt32()
	nextInstruction := state.Frame.Stack.Pop().toInt32()

	if check == 1 {
		*pc = uint(nextInstruction)
	} else {
		*pc++
	}
	return nil
}

func AllocateArrayOp(state *State, _ *GlobalData) error {

	size := state.Frame.Stack.Pop().toInt32()
	state.Frame.localVariables = make([]chan OracleConnection.Response, size)

	return nil
}

func SubscribeOp(state *State, globalData *GlobalData) error {
	stack := state.Frame.Stack
	key := stack.Pop()
	size := stack.Pop()
	offset := stack.Pop()
	index := stack.Pop()
	url, err := state.Memory.loadString(int(offset.toInt32()), int(size.toInt32()))
	if err != nil {
		return err
	}

	state.Frame.localVariables[index.toInt32()] = make(chan OracleConnection.Response)
	sub := &OracleConnection.SubscribeMessage{
		OracleKey: key.toString(),
		Url:       url,
		ResChan:   state.Frame.localVariables[index.toInt32()],
	}
	globalData.OraclePool.SubscribeChannel <- sub
	state.Frame.buffer = append(state.Frame.buffer, sub)
	return nil
}
