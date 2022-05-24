package VM

import (
	"github.com/mahmednabil109/gdeb/OracleConnection"
)

type (
	OperationType func(interpreter *Interpreter) error
)

func AddOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.Add(b))

	return nil
}

func SubOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.Sub(b))

	return nil
}

func MulOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	result, _ := a.Multiply(b)
	stack.Push(result)

	return nil
}

//GreaterOp Return 1 if a > b
func GreaterOp(interpreter *Interpreter) error {
	state := interpreter.state
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

func XorOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.Xor(b))

	return nil
}
func AndOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.And(b))

	return nil
}
func OrOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Frame.Stack

	b := stack.Pop()
	a := stack.Pop()

	stack.Push(a.Or(b))
	return nil
}

func NotOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Frame.Stack

	a := stack.Pop()

	stack.Push(a.Not())
	return nil
}

func PushOp(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Frame.Stack
	newData := NewDataWord()
	newData.SetDataWord(interpreter.ContractCode[state.Frame.pc+1 : state.Frame.pc+33])
	stack.Push(newData)
	return nil

}

func PopOp(interpreter *Interpreter) error {
	state := interpreter.state
	state.Frame.Stack.Pop()
	return nil

}

func MStoreOp(interpreter *Interpreter) error {
	state := interpreter.state
	mem := state.Memory

	offset, value := state.Frame.Stack.Pop().toInt32(), state.Frame.Stack.Pop().toByteArray()
	mem.store(int(offset), value)
	return nil

}

func MLoadOp(interpreter *Interpreter) error {
	state := interpreter.state
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

func JumpOp(interpreter *Interpreter) error {
	state := interpreter.state
	pc := &state.Frame.pc
	*pc = uint(state.Frame.Stack.Pop().toInt32())
	return nil
}

// JumpIOp conditional Jump
func JumpIOp(interpreter *Interpreter) error {
	state := interpreter.state
	pc := &state.Frame.pc
	nextInstruction := state.Frame.Stack.Pop().toInt32()
	check := state.Frame.Stack.Pop().toInt32()

	if check == 1 {
		*pc = uint(nextInstruction)
	} else {
		*pc++
	}
	return nil
}

func AllocateArrayOp(interpreter *Interpreter) error {
	state := interpreter.state

	size := state.Frame.Stack.Pop().toInt32()
	state.Frame.localVariables = make([]OracleConnection.BroadcastMsg, size)

	return nil
}

func SubscribeOp(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Frame.Stack
	localVariableIndex := stack.Pop()
	keyType := stack.Pop()
	key := stack.Pop()
	size := stack.Pop()
	offset := stack.Pop()

	url, err := state.Memory.loadString(int(offset.toInt32()), int(size.toInt32()))
	if err != nil {
		return err
	}

	sub := &OracleConnection.SubscribeMsg{
		VmId:          interpreter.Id,
		OracleKey:     key.toString(),
		KeyType:       int(keyType.toInt32()),
		Url:           url,
		BroadcastChan: interpreter.ReceiveChan,
		Index:         int(localVariableIndex.toInt32()),
	}
	state.Frame.buffer = append(state.Frame.buffer, sub)
	return nil
}

// FetchDataOp TODO return error if no such oracle
func FetchDataOp(interpreter *Interpreter) error {

	return nil
}
