package VM

import (
	"errors"
	"fmt"
	"github.com/mahmednabil109/gdeb/OracleListener"
)

type Interpreter struct {
	Id                 int
	state              *State
	ContractCode       []byte
	OraclePool         *OracleListener.OraclePool
	ReceiveChan        chan *OracleListener.BroadcastMsg
	jumpTable          *JumpTable
	gasLimit           uint64
	oracleTransactions map[int]*OracleListener.BroadcastMsg
}

func newInterpreter(id int, contractCode []byte, pool *OracleListener.OraclePool, gasLimit uint64) *Interpreter {
	return &Interpreter{
		Id:                 id,
		state:              newState(),
		ContractCode:       contractCode,
		OraclePool:         pool,
		ReceiveChan:        make(chan *OracleListener.BroadcastMsg),
		jumpTable:          newJumpTable(),
		gasLimit:           gasLimit,
		oracleTransactions: make(map[int]*OracleListener.BroadcastMsg),
	}
}

// TODO when error happens --> unsubscribe from the oraclePool
func (interpreter *Interpreter) execute() error {

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

		*consumedGas += operationInfo.gasPrice

		*pc += operationInfo.pcJump

		fmt.Println(interpreter.state)
		if int(*pc) >= len(interpreter.ContractCode) {
			return errors.New("cannot reach the globalData")
		}
	}
}
