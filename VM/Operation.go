package VM

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/mahmednabil109/gdeb/Listeners/TimeListener"
	"github.com/mahmednabil109/gdeb/Messages"
	"github.com/mahmednabil109/gdeb/VM/DataTypes"
	"github.com/mahmednabil109/gdeb/data"
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
		Data:     c[0:len(a.Data)],
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

func DivOP(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	b := stack.Pop()
	a := stack.Pop()

	err := checkCompatability("DIV", a, b)
	if err != nil {
		return err
	}
	c, err := a.Data.Div(b.Data)
	if err != nil {
		return err
	}
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

	result := DataTypes.NewData(DataTypes.Boolean)

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

	result := DataTypes.NewData(DataTypes.Boolean)

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

	result := DataTypes.NewData(DataTypes.Boolean)

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

	result := DataTypes.NewData(DataTypes.Boolean)

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

	result := DataTypes.NewData(DataTypes.Boolean)

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

	result := DataTypes.NewData(DataTypes.Boolean)

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

	result := DataTypes.NewData(DataTypes.Boolean)

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
	if a.Datatype == DataTypes.Boolean {
		c = DataTypes.NewExtraBigInt(1)
		if a.Data[0] == 0 {
			c[0] = 1
		} else {
			c[0] = 0
		}
	}
	result := &DataTypes.DataWord{
		Data:     c,
		Datatype: a.Datatype,
	}
	stack.Push(result)
	return nil
}
func IsZeroOp(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	a := stack.Pop()
	if a.Datatype == DataTypes.String || a.Datatype == DataTypes.Time {
		return errors.New("cannot perform not on " + strconv.Itoa(int(a.Datatype)))
	}
	c := a.Data.IsZero()
	d := DataTypes.NewData(DataTypes.Boolean)
	if c == true {
		d.Data[0] = 1
	}
	result := d
	stack.Push(result)
	return nil
}

func IsNegativeOp(interpreter *Interpreter) error {
	state := interpreter.state
	stack := state.Stack

	a := stack.Pop()
	if a.Datatype == DataTypes.String || a.Datatype == DataTypes.Time {
		return errors.New("cannot perform not on " + strconv.Itoa(int(a.Datatype)))
	}
	c := a.Data.Sign()
	d := DataTypes.NewData(DataTypes.Boolean)
	if c == 1 {
		d.Data[0] = 1
	}
	result := d
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
	result := DataTypes.NewData(datatype)
	result.Data.SetDataWord(code[state.pc+1 : state.pc+argSize+1])

	*pc += uint(argSize) + 1
	stack.Push(result)
	return nil
}

func PushStringOp(interpreter *Interpreter) error {
	stack := interpreter.state.Stack
	size := stack.Pop().Data.ToInt32()
	offset := stack.Pop().Data.ToInt32()

	s := interpreter.ContractCode[offset : offset+size]

	result := DataTypes.NewString(s)
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

func AddToPCOp(interpreter *Interpreter) error {
	state := interpreter.state
	pc := &state.pc
	*pc += uint(state.Stack.Pop().Data.ToInt32())
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
	keySize := stack.Pop().Data.ToInt32()
	keyOffset := stack.Pop().Data.ToInt32()
	urlSize := int(stack.Pop().Data.ToInt32())
	urlOffset := int(stack.Pop().Data.ToInt32())
	url := string(interpreter.ContractCode[urlOffset : urlOffset+urlSize])
	key := string(interpreter.ContractCode[keyOffset : keyOffset+keySize])
	sub := &Messages.SubscribeMsg{
		VmId:          interpreter.Id,
		OracleKey:     key,
		KeyType:       int(keyType.ToInt32()),
		Url:           url,
		BroadcastChan: interpreter.ReceiveChan,
		Index:         int(returnIndex),
	}
	interpreter.reservedIndex = append(interpreter.reservedIndex, int(returnIndex))
	Pool.Subscribe(sub)

	return nil
}

func WaitOp(interpreter *Interpreter) error {
	receiveChan := interpreter.ReceiveChan
	Memory := interpreter.state.Memory

Outer:
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
			d := DataTypes.NewData(DataTypes.Datatype(dataType))
			d.Data.SetDataWord(val)
			log.Println("VM Received update on", msg.Key, "=", d.Data.ToInt32())
			fmt.Println("VM Received update on", msg.Key, "=", d.Data.ToInt32())
			Memory[idx] = d
			for i, v := range interpreter.reservedIndex {
				if Memory[v] == nil {
					fmt.Println(i, "is nil")
					continue Outer
				}
			}
			if !interpreter.IsBlocked {
				log.Println("Start Evaluation")
				fmt.Println("Start Evaluation")
				return nil
			}
		case status := <-interpreter.getStatusChan:
			interpreter.IsBlocked = status
			for _, v := range interpreter.reservedIndex {
				if Memory[v] == nil {
					fmt.Println("Start Evaluation")
					continue Outer
				}
			}
			if !interpreter.IsBlocked {
				log.Println("Start Evaluation")
				fmt.Println("Start Evaluation")
				return nil
			}
		}

	}
}

func LoadOp(interpreter *Interpreter) error {
	memory, idx := interpreter.state.Memory, interpreter.state.Stack.Pop().Data.ToInt32()
	interpreter.state.Stack.Push(memory[idx])
	return nil
}

func StoreOp(interpreter *Interpreter) error {
	memory, idx, word := interpreter.state.Memory, interpreter.state.Stack.Pop().Data.ToInt32(), interpreter.state.Stack.Pop()
	memory[idx] = word
	return nil
}

func ExecutePeriodicallyOp(interpreter *Interpreter) error {

	stack := interpreter.state.Stack
	executionInterval := stack.Pop()
	freq := stack.Pop().Data.ToInt32()
	startTime := stack.Pop()
	if freq < uint32(TimeListener.ONCE) || freq > uint32(TimeListener.EveryMinute) {
		return errors.New("error in Frequency")
	}
	if executionInterval.Datatype == DataTypes.Boolean || executionInterval.Datatype == DataTypes.Time ||
		executionInterval.Datatype == DataTypes.String {
		return errors.New("cannot use " + (strconv.Itoa(int(executionInterval.Datatype))) + " as ExecutionInterval")
	}
	msg := &TimeListener.SubscribeMsg{
		Id:            interpreter.Id,
		TimeArr:       startTime.Data.ToByteArray(),
		Frequency:     TimeListener.Frequency(freq),
		Interval:      time.Duration(executionInterval.Data.ToInt32()) * time.Minute,
		GetStatusChan: interpreter.getStatusChan,
	}
	interpreter.TimeListener.Subscribe(msg)
	interpreter.IsBlocked = true
	interpreter.IsTimeDependent = true
	return nil
}

// Transfer TODO
func Transfer(interpreter *Interpreter) error {
	stack := interpreter.state.Stack
	money := stack.Pop()
	fmt.Println()
	b := stack.Pop().Data.ToString()
	a := stack.Pop().Data.ToString()
	fmt.Println(a)
	fmt.Println(b)

	// Add ConsumedGas to Transaction as field in all cases (crashes or completes)
	// reset ConsumedGas variable when new transaction sent to the channel
	// keep track of global gas limit
	//interpreter.TransactionChan <- struct {
	//	ADDRESS1 string
	//	ADDRESS2 string
	//	Money    uint64
	//  ConsumedGas unit64
	//}{
	//	ADDRESS1: add1,
	//	ADDRESS2: add2,
	//	Money:    money.Data.ToInt64(),
	//}
	interpreter.TransactionChan <- data.Transaction{
		Amount:      money.Data.ToInt64(),
		From:        a,
		To:          b,
		ConsumedGas: interpreter.state.consumedGas,
		// Timestamp:   time.Now().String(),
	}

	interpreter.state.consumedGas = 0
	log.Println("Transfer", money.Data.ToInt64(), "from", a, " to ", b)
	fmt.Println("Transfer", money.Data.ToInt64(), "from", a, " to ", b)
	if interpreter.IsTimeDependent {
		interpreter.IsBlocked = true
	}
	return nil
}
