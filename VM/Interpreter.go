package VM

import (
	"errors"
	"fmt"
	"github.com/mahmednabil109/gdeb/OracleConnection"
)

type GlobalData struct {
	ContractCode []byte
	OraclePool   OracleConnection.OraclePool
}

type Interpreter struct {
	state            *State
	globalData       *GlobalData
	operationMapping *OperationMapping
	gasLimit         uint64
}

func newInterpreter(code *GlobalData, gasLimit uint64) *Interpreter {
	return &Interpreter{state: newVM(), globalData: code, operationMapping: newInstructionInfo(), gasLimit: gasLimit}
}

func (interpreter *Interpreter) execute() error {

	for {

		consumedGas, pc := &interpreter.state.consumedGas, &interpreter.state.Frame.pc
		gasLimit := interpreter.gasLimit

		if *consumedGas > gasLimit {
			return errors.New("out of gas error")
		}

		curInstruction := (interpreter.globalData.ContractCode)[*pc]
		operationInfo := interpreter.operationMapping.getInstruction(curInstruction)

		if interpreter.state.Frame.Stack.Size() < operationInfo.stackArgsCount {
			//stack underflow exception
			return errors.New("stack underflow error")
		}

		err := operationInfo.execute(interpreter.state, interpreter.globalData)
		if err != nil {
			return err
		}

		*consumedGas += operationInfo.gasPrice

		*pc += operationInfo.pcJump

		fmt.Println(interpreter.state.toString())
		if int(*pc) >= len(interpreter.globalData.ContractCode) {
			return errors.New("cannot reach the globalData")
		}
	}
}
