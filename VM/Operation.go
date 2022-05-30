package VM

import (
	"errors"
	"github.com/mahmednabil109/gdeb/OracleListener"
	"github.com/mahmednabil109/gdeb/VM/DataTypes"
	"strconv"
)

type (
	OperationType func(interpreter *Interpreter) error
)

func checkCompatability(operation string, a, b *DataTypes.DataWord) error {

	if a.Datatype != b.Datatype {
		return errors.New("incompatible datatypes" + strconv.Itoa(int(a.Datatype)) + " and " + strconv.Itoa(int(b.Datatype)))
	}

	switch operation {
	case "ADD", "MUL", "DIV", "SUB":
		if a.Datatype == DataTypes.Time || a.Datatype == DataTypes.String || a.Datatype == DataTypes.Boolean {
			return errors.New("cannot perform" + operation +
				" operation on type" + strconv.Itoa(int(a.Datatype)))
		}
		if b.Datatype == DataTypes.Time || b.Datatype == DataTypes.String || b.Datatype == DataTypes.Boolean {
			return errors.New("cannot perform " + operation + " operation on type" + strconv.Itoa(int(a.Datatype)))
		}
	case "GT", "LT", "SLT", "SGT", "GTEQ", "LTEQ":
		if a.Datatype == DataTypes.Boolean {
			return errors.New("cannot perform" + operation +
				" operation on type" + strconv.Itoa(int(a.Datatype)))
		}
		if b.Datatype == DataTypes.Boolean {
			return errors.New("cannot perform " + operation + " operation on type" + strconv.Itoa(int(a.Datatype)))
		}
	case "XOR", "AND", "OR", "NOT":
		if a.Datatype == DataTypes.Time || a.Datatype == DataTypes.String {
			return errors.New("cannot perform" + operation +
				" operation on type" + strconv.Itoa(int(a.Datatype)))
		}
		if b.Datatype == DataTypes.Time || b.Datatype == DataTypes.String {
			return errors.New("cannot perform " + operation + " operation on type" + strconv.Itoa(int(a.Datatype)))
		}

	}

	return nil
}

func AddOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("ADD", a, b)
	if err != nil {
		return err
	}
	c := a.Data.Add(b.Data)
	result := &DataTypes.DataWord{
		Data:     c,
		Datatype: a.Datatype,
	}
	stack.Push(result)

	return nil
}

func SubOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("SUB", a, b)
	if err != nil {
		return err
	}
	c := a.Data.Sub(b.Data)
	result := &DataTypes.DataWord{
		Data:     c,
		Datatype: a.Datatype,
	}
	stack.Push(result)

	return nil
}

func MulOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("MUL", a, b)
	if err != nil {
		return err
	}
	c, _ := a.Data.Multiply(b.Data)
	result := &DataTypes.DataWord{
		Data:     c,
		Datatype: a.Datatype,
	}
	stack.Push(result)

	return nil
}

//GreaterOp Return 1 if a > b
func GreaterOp(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("GT", a, b)
	if err != nil {
		return err
	}
	isGreater := a.Data.GT(b.Data)

	var result *DataTypes.DataWord

	if isGreater {
		result.Data[0] = 1
	} else {
		result.Data[0] = 0
	}

	result.Datatype = DataTypes.Boolean

	stack.Push(result)

	return nil
}

func LessOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("LT", a, b)
	if err != nil {
		return err
	}
	isLess := a.Data.LT(b.Data)

	var result *DataTypes.DataWord

	if isLess {
		result.Data[0] = 1
	} else {
		result.Data[0] = 0
	}

	result.Datatype = DataTypes.Boolean

	stack.Push(result)

	return nil
}

func SGreaterOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("SGT", a, b)
	if err != nil {
		return err
	}
	isGreater := a.Data.SGT(b.Data)

	var result *DataTypes.DataWord

	if isGreater {
		result.Data[0] = 1
	} else {
		result.Data[0] = 0
	}

	result.Datatype = DataTypes.Boolean

	stack.Push(result)

	return nil
}

func SLessOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("SLT", a, b)
	if err != nil {
		return err
	}
	isLess := a.Data.SLT(b.Data)

	var result *DataTypes.DataWord

	if isLess {
		result.Data[0] = 1
	} else {
		result.Data[0] = 0
	}

	result.Datatype = DataTypes.Boolean

	stack.Push(result)

	return nil
}

func GTEQ_OP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("GTEQ", a, b)
	if err != nil {
		return err
	}
	isGTEQ := a.Data.GT(b.Data) || a.Data.Eq(b.Data)

	var result *DataTypes.DataWord

	if isGTEQ {
		result.Data[0] = 1
	} else {
		result.Data[0] = 0
	}

	result.Datatype = DataTypes.Boolean

	stack.Push(result)

	return nil
}

func LTEQ_OP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("GTEQ", a, b)
	if err != nil {
		return err
	}
	isLTEQ := a.Data.LT(b.Data) || a.Data.Eq(b.Data)

	var result *DataTypes.DataWord

	if isLTEQ {
		result.Data[0] = 1
	} else {
		result.Data[0] = 0
	}

	result.Datatype = DataTypes.Boolean

	stack.Push(result)

	return nil
}

func EqOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("SLT", a, b)
	if err != nil {
		return err
	}
	isEQ := a.Data.Eq(b.Data)

	var result *DataTypes.DataWord

	if isEQ {
		result.Data[0] = 1
	} else {
		result.Data[0] = 0
	}

	result.Datatype = DataTypes.Boolean

	stack.Push(result)

	return nil
}

func XorOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("XOR", a, b)
	if err != nil {
		return err
	}
	c := a.Data.Xor(b.Data)
	result := &DataTypes.DataWord{
		Data:     c,
		Datatype: a.Datatype,
	}
	stack.Push(result)

	return nil

}
func AndOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("AND", a, b)
	if err != nil {
		return err
	}
	c := a.Data.And(b.Data)
	result := &DataTypes.DataWord{
		Data:     c,
		Datatype: a.Datatype,
	}
	stack.Push(result)

	return nil
}
func OrOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("OR", a, b)
	if err != nil {
		return err
	}
	c := a.Data.Or(b.Data)
	result := &DataTypes.DataWord{
		Data:     c,
		Datatype: a.Datatype,
	}
	stack.Push(result)

	return nil
}

func NotOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	a := stack.Pop()
	if a.Datatype == DataTypes.String || a.Datatype == DataTypes.Time {
		return errors.New("cannot perform not on " + strconv.Itoa(int(a.Datatype)))
	}
	c := a.Data.Not()
	result := &DataTypes.DataWord{
		Data:     c,
		Datatype: a.Datatype,
	}
	stack.Push(result)
	return nil
}

func PushOp(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack
	opCode := interpreter.ContractCode[state.pc]
	jumpTable := interpreter.jumpTable
	code := interpreter.ContractCode
	pc := &state.pc

	var result *DataTypes.DataWord
	var datatype DataTypes.Datatype
	var argSize uint

	switch opCode {
	case byte(PUSH32):
		argSize = jumpTable[PUSH32].codeArgsCount
		datatype = DataTypes.Int32
	case byte(PUSHBOOl):
		argSize = jumpTable[PUSHBOOl].codeArgsCount
		datatype = DataTypes.Boolean
	case byte(PUSHTIME):
		argSize = jumpTable[PUSHTIME].codeArgsCount
		datatype = DataTypes.Time
	case byte(PUSH64):
		argSize = jumpTable[PUSH64].codeArgsCount
		datatype = DataTypes.Int64
	case byte(PUSH256):
		argSize = jumpTable[PUSH256].codeArgsCount
		datatype = DataTypes.Int256
	}
	result.Data.SetDataWord(code[state.pc+1 : state.pc+argSize+1])
	result.Datatype = datatype
	*pc += uint(argSize) + 1

	stack.Push(result)
	return nil
}

func PopOp(interpreter *Interpreter) error {
	state := interpreter.state
	state.Stack.Pop()
	return nil
}

func JumpOp(interpreter *Interpreter) error {
	state := interpreter.state
	pc := &state.pc
	*pc = uint(state.Stack.Pop().Data.ToInt32())
	return nil
}

// JumpIOp conditional Jump
func JumpIOp(interpreter *Interpreter) error {
	state := interpreter.state
	pc := &state.pc
	nextInstruction := state.Stack.Pop().Data.ToInt32()
	check := state.Stack.Pop().Data.ToInt32()

	if check == 1 {
		*pc = uint(nextInstruction)
	} else {
		*pc++
	}
	return nil
}

func SubscribeOp(interpreter *Interpreter) error {
	Pool := interpreter.OraclePool
	state := interpreter.state
	stack := state.Stack
	returnIndex := stack.Pop().Data.ToInt32()
	keyType := stack.Pop().Data
	key := stack.Pop().Data
	size := stack.Pop().Data
	offset := stack.Pop().Data

	url := string(interpreter.ContractCode[int(offset.ToInt32()):int(size.ToInt32())])

	sub := &OracleListener.SubscribeMsg{
		VmId:          interpreter.Id,
		OracleKey:     string(key.ToInt32()),
		KeyType:       int(keyType.ToInt32()),
		Url:           url,
		BroadcastChan: interpreter.ReceiveChan,
		Index:         int(returnIndex),
	}
	Pool.Subscribe(sub)

	return nil
}

func FetchDataOp(interpreter *Interpreter) error {
	receiveChan := interpreter.ReceiveChan

	Memory := interpreter.state.Memory
	for {

		select {
		case msg := <-receiveChan:
			if msg.Error {
				return errors.New("no such oracle")
			}
			idx := msg.Index
			val := msg.Value
			dataType := msg.Type
			interpreter.oracleTransactions[idx] = msg
			Memory[idx] = &DataTypes.DataWord{
				Data:     DataTypes.ByteArrToBigInt(val),
				Datatype: DataTypes.Datatype(dataType),
			}
			for _, v := range Memory {
				if v == nil {
					break
				}
			}
			return nil
		}
	}
}

func LoadOp(interpreter *Interpreter) error {
	memory, idx := interpreter.state.Memory, interpreter.state.Stack.Pop().Data.ToInt32()
	interpreter.state.Stack.Push(memory[idx])
	return nil
}
