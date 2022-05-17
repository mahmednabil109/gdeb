package VM

import (
	"errors"
	"fmt"
)

type ContractByteCode []byte

type Interpreter struct {
	state            *VMState
	code             *ContractByteCode
	operationMapping *OperationMapping
	gasLimit         uint64
}

func newInterpreter(code *ContractByteCode, gasLimit uint64) *Interpreter {
	return &Interpreter{state: newVM(), code: code, operationMapping: newInstructionInfo(), gasLimit: gasLimit}
}

func (interpreter *Interpreter) execute() error {

	for {

		consumedGas, pc := &interpreter.state.consumedGas, &interpreter.state.Frame.pc
		gasLimit := interpreter.gasLimit

		if *consumedGas > gasLimit {
			return errors.New("out of gas error")
		}

		curInstruction := (*interpreter.code)[*pc]
		operationInfo := interpreter.operationMapping.getInstruction(curInstruction)

		if interpreter.state.Frame.Stack.Size() < operationInfo.stackArgsCount {
			//stack underflow exception
			return errors.New("stack underflow error")
		}

		err := operationInfo.execute(interpreter.state, interpreter.code)
		if err != nil {
			return err
		}

		*consumedGas += operationInfo.gasPrice

		*pc += operationInfo.pcJump

		fmt.Println(interpreter.state.toString())
		if int(*pc) >= len(*interpreter.code) {
			return errors.New("cannot reach the code")
		}
	}
}
