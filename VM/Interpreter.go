package VM

import (
	"errors"
	"fmt"
	"github.com/mahmednabil109/gdeb/Listeners/OracleListener"
	"github.com/mahmednabil109/gdeb/Messages"
	"time"
)

type Frequency byte

const (
	ONCE    Frequency = 0
	HOURLY  Frequency = 1
	DAILY   Frequency = 2
	WEEKLY  Frequency = 3
	MONTHLY Frequency = 4
	YEARLY  Frequency = 5
)

type PeriodicExecution struct {
	startTime         time.Time
	frequency         Frequency
	executionInterval time.Duration // in minutes
}

type Interpreter struct {
	Id                 int
	state              *State
	ContractCode       []byte
	OraclePool         *OracleListener.OraclePool
	ReceiveChan        chan *Messages.BroadcastMsg
	TransactionChan    chan interface{}
	jumpTable          *JumpTable
	gasLimit           uint64
	oracleTransactions map[int]*Messages.BroadcastMsg
	reservedIndex      []int
	periodicExecution  *PeriodicExecution
}

func newInterpreter(id int, contractCode []byte, pool *OracleListener.OraclePool, gasLimit uint64) *Interpreter {
	return &Interpreter{
		Id:                 id,
		state:              newState(),
		ContractCode:       contractCode,
		OraclePool:         pool,
		ReceiveChan:        make(chan *Messages.BroadcastMsg),
		jumpTable:          newJumpTable(),
		gasLimit:           gasLimit,
		oracleTransactions: make(map[int]*Messages.BroadcastMsg),
		reservedIndex:      make([]int, 0),
	}
}

// TODO when error happens --> unsubscribe from the oraclePool
func (interpreter *Interpreter) run() error {

	for {

		consumedGas, pc := &interpreter.state.consumedGas, &interpreter.state.pc
		gasLimit := interpreter.gasLimit

		if *consumedGas > gasLimit {
			return errors.New("out of gas error")
		}

		curInstruction := (interpreter.ContractCode)[*pc]
		operationInfo := interpreter.jumpTable.getInstruction(curInstruction)

		if interpreter.state.Stack.Size() < operationInfo.stackArgsCount {
			//stack underflow exception
			return errors.New("stack underflow error")
		}

		err := operationInfo.execute(interpreter)
		if err != nil {
			return err
		}

		//Todo check that the operation is not transfer
		*consumedGas += operationInfo.gasPrice

		*pc += operationInfo.pcJump

		fmt.Println(interpreter.state)
		if int(*pc) >= len(interpreter.ContractCode) {
			return errors.New("cannot reach the globalData")
		}
	}
}
