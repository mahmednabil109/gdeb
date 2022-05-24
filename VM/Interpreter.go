package VM

import (
	"errors"
	"fmt"
	"github.com/mahmednabil109/gdeb/OracleListener"
)

type GlobalData struct {
}

type Interpreter struct {
	Id               int
	state            *State
	ContractCode     []byte
	OraclePool       *OracleListener.OraclePool
	ReceiveChan      chan *OracleListener.BroadcastMsg
	operationMapping *OperationMapping
	gasLimit         uint64
}

func newInterpreter(id int, contractCode []byte, pool *OracleListener.OraclePool, gasLimit uint64) *Interpreter {
	return &Interpreter{
		Id:               id,
		state:            newVM(),
		ContractCode:     contractCode,
		OraclePool:       pool,
		ReceiveChan:      make(chan *OracleListener.BroadcastMsg),
		operationMapping: newInstructionInfo(),
		gasLimit:         gasLimit,
	}
}

// TODO when error happens --> unsubscribe from the oraclePool
func (interpreter *Interpreter) execute() error {

	for {

		consumedGas, pc := &interpreter.state.consumedGas, &interpreter.state.Frame.pc
		gasLimit := interpreter.gasLimit

		if *consumedGas > gasLimit {
			return errors.New("out of gas error")
		}

		curInstruction := (interpreter.ContractCode)[*pc]
		operationInfo := interpreter.operationMapping.getInstruction(curInstruction)

		if interpreter.state.Frame.Stack.Size() < operationInfo.stackArgsCount {
			//stack underflow exception
			return errors.New("stack underflow error")
		}

		err := operationInfo.execute(interpreter)
		if err != nil {
			return err
		}

		*consumedGas += operationInfo.gasPrice

		*pc += operationInfo.pcJump

		fmt.Println(interpreter.state.toString())
		if int(*pc) >= len(interpreter.ContractCode) {
			return errors.New("cannot reach the globalData")
		}
	}
}
